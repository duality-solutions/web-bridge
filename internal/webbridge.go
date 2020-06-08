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
	fmt.Println("Version:", version, "Hash", githash)
	fmt.Println("OS: ", runtime.GOOS)
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
	if debug {
		fmt.Println("Running WebBridge in debug log mode.")
		fmt.Println("Args", args)
	}
	if test {
		fmt.Println("Running WebBridge in test mode.")
	}
	config.Load()
	if debug {
		fmt.Println("Config", config)
	}

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
	_, err := bridge.ConnectToIceServices(config)
	if err != nil {
		return fmt.Errorf("NewPeerConnection %v", err)
	}
	if debug {
		fmt.Println("Connected to ICE services.")
	}
	proc, err := dynamic.FindDynamicdProcess()
	if err == nil {
		fmt.Println("dynamicd already running. Attempting to kill the process.")
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
		fmt.Println("Running dynamicd process found Pid", proc.Pid)
	} else {
		fmt.Println(err)
		// start again or exit app ???
	}
	if !test {
		// make sure wallet is created
		dynamicd.WaitForWalletCreated()
		status, errStatus := dynamicd.GetSyncStatus()
		if errStatus != nil {
			return fmt.Errorf("GetSyncStatus %v", errStatus)
		}
		fmt.Println("dynamicd running... Sync " + fmt.Sprintf("%f", status.SyncProgress*100) + " percent complete!")

		info, errInfo := dynamicd.GetInfo()
		if errInfo != nil {
			return fmt.Errorf("GetInfo %v", errInfo)
		}
		fmt.Println("dynamic connections", info.Connections)

		acc, errAccounts := dynamicd.GetMyAccounts()
		if errAccounts != nil {
			fmt.Println("GetActiveLinks error", errAccounts)
		} else {
			for i, account := range *acc {
				fmt.Println("Account", i+1, account.CommonName, account.ObjectFullPath, account.WalletAddress, account.LinkAddress)
			}
		}
		errUnlock := dynamicd.UnlockWallet("")
		if errUnlock != nil {
			fmt.Println("Wallet locked. Please unlock the wallet to continue.")
			var unlocked = false
			for !unlocked {
				fmt.Print("wallet passphrase> ")
				bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				walletpassphrase = strings.Trim(string(bytePassword), "\r\n ")
				errUnlock = dynamicd.UnlockWallet(walletpassphrase)
				if errUnlock == nil {
					fmt.Println("Wallet unlocked.")
					unlocked = true
				} else {
					fmt.Println(errUnlock)
				}
			}
		}
		al, errLinks := dynamicd.GetActiveLinks()
		if errLinks != nil {
			fmt.Println("GetActiveLinks error", errLinks)
		} else {
			fmt.Printf("Found %v links\n", len(al.Links))
		}
		stopWatcher := make(chan struct{})
		dynamic.WatchProcess(stopWatcher, 10, walletpassphrase)
		// TODO: REST API running
		// TODO: Admin console running
		var sync = false
		stopBridges := make(chan struct{})
		if acc != nil && al != nil {
			go bridge.StartBridges(stopBridges, config, *dynamicd, *acc, *al)
		}
		for !shutdown {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("web-bridge> ")
			cmdText, _ := reader.ReadString('\n')
			if len(cmdText) > 1 {
				cmdText = strings.Trim(cmdText, "\r\n ") //cmdText[:len(cmdText)-2]
			}
			if strings.HasPrefix(cmdText, "exit") || strings.HasPrefix(cmdText, "shutdown") || strings.HasPrefix(cmdText, "stop") {
				fmt.Println("Exit command. Stopping services.")
				shutdown = true
				close(stopWatcher)
				close(stopBridges)
				break
			} else if strings.HasPrefix(cmdText, "dynamic-cli") {
				req, errNewRequest := dynamic.NewRequest(cmdText)
				if errNewRequest != nil {
					fmt.Println("Error:", errNewRequest)
				} else {
					strResp, _ := util.BeautifyJSON(<-dynamicd.ExecCmdRequest(req))
					fmt.Println(strResp)
				}
			} else {
				fmt.Println("Invalid command", cmdText)
				status, errStatus = dynamicd.GetSyncStatus()
			}
			status, errStatus = dynamicd.GetSyncStatus()
			if errStatus != nil {
				fmt.Println("syncstatus unmarshal error", errStatus)
			} else {
				if !sync {
					fmt.Println("Sync " + fmt.Sprintf("%f", status.SyncProgress*100) + " percent complete!")
					if status.SyncProgress == 1 {
						sync = true
					}
				}
			}
		}
		bridge.ShutdownBridges()
		// Stop dynamicd
		reqStop, _ := dynamic.NewRequest("dynamic-cli stop")
		respStop, _ := util.BeautifyJSON(<-dynamicd.ExecCmdRequest(reqStop))
		fmt.Println(respStop)
		time.Sleep(time.Second * 5)
	}

	fmt.Println("Looking for dynamicd process pid", dynamicd.Cmd.Process.Pid)
	_, errFindProcess := os.FindProcess(dynamicd.Cmd.Process.Pid)
	if errFindProcess == nil {
		fmt.Println("Process found. Killing dynamicd process.")
		if errKill := dynamicd.Cmd.Process.Kill(); err != errKill {
			fmt.Println("failed to kill process: ", errKill)
		}
	} else {
		fmt.Println("Dynamicd process not found")
	}

	fmt.Println("Exit.")
	return nil
}
