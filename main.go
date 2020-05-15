package main

import (
	"fmt"
	"runtime"

	webbridge "github.com/duality-solutions/web-bridge/internal"
)

// Version is the WebBridge version number in ISO date format
var Version string

// GitHash is the git hash tag
var GitHash string

// BuildDateTime is the date the binary was built
//var BuildDateTime string

func main() {
	fmt.Println("Version:", Version, "Hash", GitHash)
	fmt.Println("OS: ", runtime.GOOS)
	webbridge.Init()
}
