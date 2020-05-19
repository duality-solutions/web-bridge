package main

import (
	"fmt"
	"os"

	webbridge "github.com/duality-solutions/web-bridge/internal"
)

// Version is the WebBridge version number in ISO date format
var Version string

// GitHash is the git hash tag
var GitHash string

// BuildDateTime is the date the binary was built
//var BuildDateTime string

func main() {
	if err := webbridge.Init(Version, GitHash); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
