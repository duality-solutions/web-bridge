package webbridge

import (
	"fmt"
	"os"
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
	// ICE service test completed

	// Dynamicd running ... sync 88% complete
	cmd, err := LoadRPCDynamicd()
	if err != nil {
		fmt.Println("Could not load dynamicd. Can not continue.", err)
		os.Exit(-1)
	}
	fmt.Println("Starting dynamicd...")
	cmd.Start()
	fmt.Println("Cmd", cmd)
	fmt.Println("Cmd Address:", &cmd)
	fmt.Printf("Cmd Type: %T\n", cmd)
	cmd.Wait()
	fmt.Println("Proc:", cmd.Process)
	fmt.Println("Proc Address:", &cmd.Process)
	fmt.Printf("Proc Type: %T\n", cmd.Process)
	// REST API running
	// Admin console running
	// Establishing WebRTC connections with links
	// Starting HTTP bridges for active links
}
