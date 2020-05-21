package webbridge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	bridge "github.com/duality-solutions/web-bridge/internal/bridge"
	dynamic "github.com/duality-solutions/web-bridge/internal/dynamic"
	settings "github.com/duality-solutions/web-bridge/internal/settings"
	util "github.com/duality-solutions/web-bridge/internal/utilities"
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
var development = false
var debug = false
var shutdown = false

// Init is used to begin all WebBridge tasks
func Init(version, githash string) error {
	fmt.Println("Version:", version, "Hash", githash)
	fmt.Println("OS: ", runtime.GOOS)
	args := os.Args[1:]
	if len(args) > 0 {
		for _, v := range args {
			switch v {
			case "-dev":
				development = true
			case "-debug":
				debug = true
			}
		}
	}
	if debug {
		fmt.Println("Running WebBridge in debug log mode.")
		fmt.Println("Args", args)
	}
	if development {
		fmt.Println("Running WebBridge in development mode.")
	}
	config.Load()
	if debug {
		fmt.Println("Config", config)
	}
	// Connect to ICE services
	peerConnection, err := bridge.ConnectToIceServices(config)
	if err != nil {
		return fmt.Errorf("NewPeerConnection", err)
	}
	if debug {
		fmt.Println("Connected to ICE services.")
	}
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return fmt.Errorf("CreateOffer", err)
	}
	if debug {
		fmt.Println(offer, "\nCreated WebRTC offer successfully!")
	}
	// create and run dynamicd
	dynamicd, err := dynamic.LoadRPCDynamicd()
	if err != nil {
		return fmt.Errorf("LoadRPCDynamicd", err)
	}
	// TODO: check if dynamicd is already running
	// TODO: REST API running
	// TODO: Admin console running
	// TODO: Establishing WebRTC connections with links
	// TODO: Starting HTTP bridges for active links
	reqStatus, _ := dynamic.NewRequest("dynamic-cli syncstatus")
	var status dynamic.SyncStatus
	errUnmarshal := json.Unmarshal([]byte(<-dynamicd.ExecCmdRequest(reqStatus)), &status)
	if errUnmarshal != nil {
		fmt.Println("syncstatus unmarshal error", errUnmarshal)
	} else {
		fmt.Println("dynamicd running... Sync " + fmt.Sprintf("%f", status.SyncProgress*100) + " percent complete!")
	}
	errUnlock := dynamicd.UnlockWallet("")
	if errUnlock != nil {
		fmt.Println("Wallet locked. Please unlock the wallet to continue.")
		var unlocked = false
		for !unlocked {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("wallet passphrase> ")
			walletpassphrase, _ := reader.ReadString('\n')
			walletpassphrase = strings.Trim(walletpassphrase, "\r\n ")
			errUnlock = dynamicd.UnlockWallet(walletpassphrase)
			if errUnlock == nil {
				fmt.Println("Wallet unlocked.")
				unlocked = true
			} else {
				fmt.Println(errUnlock)
			}
		}
	}
	if development {
		fmt.Println("Development mode. Skipping terminal input.")
		time.Sleep(time.Second * 15)
	} else {
		for !shutdown {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("web-bridge> ")
			cmdText, _ := reader.ReadString('\n')
			if len(cmdText) > 1 {
				cmdText = strings.Trim(cmdText, "\r\n ") //cmdText[:len(cmdText)-2]
			}
			if strings.HasPrefix(cmdText, "exit") || strings.HasPrefix(cmdText, "shutdown") {
				fmt.Println("Exit command. Stopping services.")
				shutdown = true
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
				errUnmarshal = json.Unmarshal([]byte(<-dynamicd.ExecCmdRequest(reqStatus)), &status)
				if errUnmarshal != nil {
					fmt.Println("syncstatus unmarshal error", errUnmarshal)
				} else {
					fmt.Println("Sync " + fmt.Sprintf("%f", status.SyncProgress*100) + " percent complete!")
				}
			}
		}
	}
	reqStop, _ := dynamic.NewRequest("dynamic-cli stop")
	respStop, _ := util.BeautifyJSON(<-dynamicd.ExecCmdRequest(reqStop))
	fmt.Println(respStop)
	time.Sleep(time.Second * 5)
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
	fmt.Println("Good bye.")
	return nil
}
