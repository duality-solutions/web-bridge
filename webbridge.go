package main

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

func main() {
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
	if development {
		// load dynamic docker image
		dynamicd, err := LoadDockerDynamicd()
		if err != nil {
			if dynamicd != nil {
				// get percent complete from syncstatus
				fmt.Println("Dynamicd running ... sync", dynamicd, "complete")
			} else {
				fmt.Println("Dynamicd not running.")
			}
		} else {
			fmt.Println("dynamicd error:", err)
		}
	} else {
		dynamicd, err := LoadRPCDynamicd()
		if err != nil {
			if dynamicd != nil {
				// get percent complete from syncstatus
				fmt.Println("Dynamicd running ... sync", dynamicd, "complete")
			} else {
				fmt.Println("Dynamicd not running.")
			}
		} else {
			fmt.Println("dynamicd error:", err)
		}
	}
	// REST API running
	// Admin console running
	// Establishing WebRTC connections with links
	// Starting HTTP bridges for active links
}
