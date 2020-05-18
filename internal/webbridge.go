package webbridge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	dynamic "github.com/duality-solutions/web-bridge/internal/dynamic"
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

var config Configuration
var development = false
var debug = false
var shutdown = false

// Init is used to begin all WebBridge tasks
func Init() {
	if debug {
		fmt.Println("Running WebBridge in debug log mode.")
	}
	if development {
		fmt.Println("Running WebBridge in development mode.")
	}
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
		fmt.Println("Args", args)
	}
	config.load()
	if debug {
		fmt.Println("Config", config)
	}
	// TODO: ICE service test completed

	dynamicd, err := dynamic.LoadRPCDynamicd()
	if err != nil {
		fmt.Println("Could not load dynamicd. Can not continue.", err)
		os.Exit(-1)
	}
	// TODO: check if dynamicd is already running
	// TODO: REST API running
	// TODO: Admin console running
	// TODO: Establishing WebRTC connections with links
	// TODO: Starting HTTP bridges for active links
	cmdStatus := "{\"method\": \"syncstatus\", \"params\": [], \"id\": 1}"
	var status dynamic.SyncStatus
	errUnmarshal := json.Unmarshal([]byte(<-dynamicd.ExecCmd(cmdStatus)), &status)
	if errUnmarshal != nil {
		fmt.Println("cmdStatus Unmarshal error", errUnmarshal)
	} else {
		fmt.Println("dynamicd running... Sync percent (", status.SyncProgress*100, "%) complete!")
	}
	for !shutdown {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("web-bridge> ")
		cmdText, _ := reader.ReadString('\n')
		if len(cmdText) > 1 {
			cmdText = cmdText[:len(cmdText)-2]
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
			// TODO: process command here.
			fmt.Println(cmdText)
			errUnmarshal = json.Unmarshal([]byte(<-dynamicd.ExecCmd(cmdStatus)), &status)
			if errUnmarshal != nil {
				fmt.Println("cmdStatus Unmarshal error", errUnmarshal)
			} else {
				fmt.Println("Sync percent (", status.SyncProgress*100, "%) complete!")
			}
		}
	}
	cmdStop := "{\"method\": \"stop\", \"params\": [], \"id\": 2}"
	resStop, _ := util.BeautifyJSON(<-dynamicd.ExecCmd(cmdStop))
	fmt.Println(resStop)
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
}
