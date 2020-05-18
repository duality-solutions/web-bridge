package webbridge

import (
	"fmt"
	"os"
	"time"
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

	dynamicd, err := LoadRPCDynamicd()
	if err != nil {
		fmt.Println("Could not load dynamicd. Can not continue.", err)
		os.Exit(-1)
	}
	// TODO: check if dynamicd is already running
	// Success
	fmt.Println("dynamicd", dynamicd)
	fmt.Println("dynamicd address:", &dynamicd)
	fmt.Printf("dynamicd Type: %T\n", dynamicd)
	fmt.Println("dynamicd proc:", dynamicd.cmd.Process)
	fmt.Println("dynamicd address:", &dynamicd.cmd.Process)
	fmt.Printf("dynamicd proc type: %T\n", dynamicd.cmd.Process)
	fmt.Println("dynamicd cmd ctx:", dynamicd.ctx)
	fmt.Println("dynamicd cmd cancelFunc:", dynamicd.cancelFunc)
	// TODO: Create dynamicd JSON RPC controller
	// TODO: Dynamicd running ... print sync percent (like 88%) complete
	cmdSyncStatus := "{\"method\": \"syncstatus\", \"params\": [], \"id\": 1}"
	c := dynamicd.execCmd(cmdSyncStatus)
	fmt.Println("cmdSyncStatus", <-c)
	cmdLinks := "{\"method\": \"link\", \"params\": [\"complete\"], \"id\": 2}"
	c = dynamicd.execCmd(cmdLinks)
	fmt.Println("cmdLinks", <-c)
	cmdFailTest := "{\"method\": \"fail\", \"params\": [\"12345\"], \"id\": 3}"
	c = dynamicd.execCmd(cmdFailTest)
	fmt.Println("cmdFailTest", <-c)
	cmdStop := "{\"method\": \"stop\", \"params\": [], \"id\": 4}"
	c = dynamicd.execCmd(cmdStop)
	fmt.Println("cmdStop", <-c)
	time.Sleep(time.Second * 5)
	fmt.Println("Looking for dynamicd process pid", dynamicd.cmd.Process.Pid)
	_, errFindProcess := os.FindProcess(dynamicd.cmd.Process.Pid)
	if errFindProcess == nil {
		fmt.Println("Process found. Killing dynamicd process.")
		if errKill := dynamicd.cmd.Process.Kill(); err != errKill {
			fmt.Println("failed to kill process: ", errKill)
		}
	} else {
		fmt.Println("Dynamicd process not found")
	}
	// TODO: REST API running
	// TODO: Admin console running
	// TODO: Establishing WebRTC connections with links
	// TODO: Starting HTTP bridges for active links
}
