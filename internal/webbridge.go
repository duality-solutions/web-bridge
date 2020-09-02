package webbridge

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	bridge "github.com/duality-solutions/web-bridge/bridge"
	settings "github.com/duality-solutions/web-bridge/init/settings"
	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rest"
	dynamic "github.com/duality-solutions/web-bridge/rpc/dynamic"
)

/*
# Show Load Status
- Configuration loaded
- ICE service test completed
- Dynamicd running ... sync 88% complete
- REST API running
- Admin console running
- Establishing WebRTC connections with links
- Starting HTTP bridges for active links
api
- RestAPI
configs
- Config
docs
- Diagrams
init
- Main.
- call config init
- Manage channels
- Manage shutdown and cleanup
dynamicd
- Manage dynamicd and JSON RPC calls
web
- AdminConsole

WebRTCBridge

2) Load HTTP Server
	- Authentication: Use OAuth
	- Web UI admin console
		- Create accounts
		- Create links
		- Blockchain status
		- Send/receive funds
		- Start/stop/view status of link bridges connection
		- link bridges stats and logs
		- link configuration
		- link permission (out of scope for v1)
	- API Server
3) Launch Dynamic daemon RPC
	- Get all accounts
	- Get all links
	- Manage process (out of scope for v1)
	- encrypt wallet
5) Load HTTP to WebRTC bridges
	- Start HTTP to WebRTC Relay
	- Connect to all links
*/

var config settings.Configuration
var debug = false
var shutdown = false
var test = false
var walletpassphrase = ""
var testCreateOffer = false
var testWaitForOffer = false

// Init is used to begin all WebBridge tasks
func Init(version, githash string) error {
	args := os.Args[1:]
	if len(args) > 0 {
		for _, v := range args {
			switch v {
			case "-debug":
				debug = true
			case "-test":
				test = true
			case "-testCreateOffer":
				testCreateOffer = true
			case "-testWaitForOffer":
				testWaitForOffer = true
			}
		}
	}
	// initilize debug.log file
	util.InitDebugLogFile(debug)
	util.Info.Println("Version:", version, "Hash", githash)
	util.Info.Println("OS: ", runtime.GOOS)
	if debug {
		util.Info.Println("Running WebBridge in debug log mode.")
		util.Info.Println("Args", args)
	}
	if test {
		util.Info.Println("Running WebBridge in test mode.")
	}
	config.Load()
	util.Info.Println("Config", config)

	if testCreateOffer || testWaitForOffer {
		if testCreateOffer {
			bridge.TestCreateOffer()
		} else if testWaitForOffer {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("paste offer> ")
			offer, _ := reader.ReadString('\n')
			bridge.TestWaitForOffer(offer)
		}
	}

	// Connect to ICE services
	_, err := bridge.ConnectToIceServicesDetached(config)
	if err != nil {
		return fmt.Errorf("ConnectToIceServicesDetached error %v", err)
	}
	util.Info.Println("Connected to ICE services.")

	proc, err := dynamic.FindDynamicdProcess()
	if err == nil {
		util.Warning.Println("dynamicd already running. Attempting to kill the process.")
		err = proc.Kill()
		if err != nil {
			return fmt.Errorf("Fatal error, dynamicd process (%v) is running but can't be stopped %v", proc.Pid, err)
		}
	}
	// create and run dynamicd
	dynamicd, err := dynamic.LoadRPCDynamicd()
	if err != nil {
		return fmt.Errorf("LoadRPCDynamicd %v", err)
	}

	proc, err = dynamic.FindDynamicdProcess()
	if proc != nil {
		util.Info.Println("Running dynamicd process found Pid", proc.Pid)
	} else {
		util.Error.Println(err)
		// start again or exit app ???
	}
	if !test {
		// make sure wallet is created
		dynamicd.WaitForWalletCreated()
		status, errStatus := dynamicd.GetSyncStatus()
		if errStatus != nil {
			return fmt.Errorf("GetSyncStatus %v", errStatus)
		}
		util.Info.Println("dynamicd running... Sync " + fmt.Sprintf("%f", status.SyncProgress*100) + " percent complete!")

		info, errInfo := dynamicd.GetInfo()
		if errInfo != nil {
			return fmt.Errorf("GetInfo %v", errInfo)
		}
		util.Info.Println("dynamic connections", info.Connections)

		acc, errAccounts := dynamicd.GetMyAccounts()
		if errAccounts != nil {
			util.Error.Println("GetActiveLinks error", errAccounts)
		} else {
			for i, account := range *acc {
				util.Info.Println("Account", i+1, account.CommonName, account.ObjectFullPath, account.WalletAddress, account.LinkAddress)
			}
		}
		var mode string = "release"
		if debug {
			mode = "debug"
		}
		// Start Gin web services
		go rest.StartWebServiceRouter(dynamicd, mode)

		errUnlock := dynamicd.UnlockWallet("")
		if errUnlock != nil {
			util.Info.Println("Wallet locked. Please unlock the wallet to continue.")
			var unlocked = false
			for !unlocked {
				fmt.Print("wallet passphrase> ")
				bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				walletpassphrase = strings.Trim(string(bytePassword), "\r\n ")
				errUnlock = dynamicd.UnlockWallet(walletpassphrase)
				if errUnlock == nil {
					util.Info.Println("Wallet unlocked.")
					unlocked = true
				} else {
					util.Error.Println(errUnlock)
				}
			}
		}
		al, errLinks := dynamicd.GetActiveLinks()
		if errLinks != nil {
			util.Error.Println("GetActiveLinks error", errLinks)
		} else {
			util.Info.Printf("Found %v links\n", len(al.Links))
		}
		stopWatcher := make(chan struct{})
		dynamic.WatchProcess(stopWatcher, 10, walletpassphrase)
		// TODO: Admin console running
		var sync = false
		stopBridges := make(chan struct{})
		if acc != nil && al != nil {
			go bridge.StartBridges(&stopBridges, config, *dynamicd, *acc, *al)
		}
		for !shutdown {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("web-bridge> ")
			cmdText, _ := reader.ReadString('\n')
			if len(cmdText) > 1 {
				cmdText = strings.Trim(cmdText, "\r\n ") //cmdText[:len(cmdText)-2]
			}
			if strings.HasPrefix(cmdText, "exit") || strings.HasPrefix(cmdText, "shutdown") || strings.HasPrefix(cmdText, "stop") {
				util.Info.Println("Exit command. Stopping services.")
				shutdown = true
				close(stopWatcher)
				break
			} else if strings.HasPrefix(cmdText, "dynamic-cli") {
				req, errNewRequest := dynamic.NewRequest(cmdText)
				if errNewRequest != nil {
					util.Error.Println("Error:", errNewRequest)
				} else {
					strResp, _ := util.BeautifyJSON(<-dynamicd.ExecCmdRequest(req))
					util.Info.Println(strResp)
				}
			} else {
				util.Warning.Println("Invalid command", cmdText)
				status, errStatus = dynamicd.GetSyncStatus()
			}
			status, errStatus = dynamicd.GetSyncStatus()
			if errStatus != nil {
				util.Error.Println("syncstatus unmarshal error", errStatus)
			} else {
				if !sync {
					util.Info.Println("Sync " + fmt.Sprintf("%f", status.SyncProgress*100) + " percent complete!")
					if status.SyncProgress == 1 {
						sync = true
					}
				}
			}
		}
		bridge.ShutdownBridges(&stopBridges)
		// Stop dynamicd
		reqStop, _ := dynamic.NewRequest("dynamic-cli stop")
		respStop, _ := util.BeautifyJSON(<-dynamicd.ExecCmdRequest(reqStop))
		util.Info.Println(respStop)
		time.Sleep(time.Second * 5)
	}

	util.Info.Println("Looking for dynamicd process pid", dynamicd.Cmd.Process.Pid)
	_, errFindProcess := os.FindProcess(dynamicd.Cmd.Process.Pid)
	if errFindProcess == nil {
		util.Info.Println("Process found. Killing dynamicd process.")
		if errKill := dynamicd.Cmd.Process.Kill(); err != errKill {
			util.Error.Println("failed to kill process: ", errKill)
		}
	} else {
		util.Info.Println("Dynamicd process not found")
	}

	util.Info.Println("Exit.")
	util.EndDebugLogFile(30)
	return nil
}
