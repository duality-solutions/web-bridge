package webbridge

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/duality-solutions/web-bridge/bridge"
	"github.com/duality-solutions/web-bridge/init/settings"
	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/duality-solutions/web-bridge/rest"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
)

const (
	// DefaultName application name
	DefaultName string = "web-bridge"
)

var config settings.Configuration
var debug = false
var test = false
var walletpassphrase = ""
var testCreateOffer = false
var testWaitForOffer = false

// Init is used to begin all WebBridge tasks
func Init(version, githash string) error {
	running, pid, err := util.FindWebBridgeProcess(DefaultName)
	if running {
		if err == nil && pid > 0 {
			return fmt.Errorf("web-bridge process (%v) found running. Can only run one web-bridge instance at a time", pid)
		}
		return fmt.Errorf("web-bridge process (%v) found running. Can only run one web-bridge instance at a time. Error: %v", pid, err)
	}
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
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	pathSeperator := ""
	if runtime.GOOS == "windows" {
		pathSeperator = `\\`
		homeDir += pathSeperator + `.` + DefaultName + pathSeperator
	} else {
		pathSeperator = `/`
		homeDir += pathSeperator + `.` + DefaultName + pathSeperator
	}
	// initilize debug.log file
	util.InitDebugLogFile(debug, homeDir)
	util.Info.Println("Version:", version, "Hash", githash)
	util.Info.Println("OS: ", runtime.GOOS)
	if debug {
		util.Info.Println("Running WebBridge in debug log mode.")
		util.Info.Println("Args", args)
	}
	if test {
		util.Info.Println("Running WebBridge in test mode.")
	}
	config.Load(homeDir, pathSeperator)
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
	_, err = bridge.ConnectToIceServicesDetached(&config)
	if err != nil {
		return fmt.Errorf("ConnectToIceServicesDetached error %v", err)
	}
	util.Info.Println("Connected to ICE services:")
	if !test {
		dynamicd, proc, err := dynamic.FindDynamicdProcess(false, time.Second*1)
		if err == nil {
			// kill existing dynamicd process
			util.Warning.Println("dynamicd daemon already running. Attempting to kill the process.")
			err = proc.Kill()
			if err != nil {
				return fmt.Errorf("Fatal error, dynamicd daemon process (%v) is running but can't be stopped %v", proc.Pid, err)
			}
			time.Sleep(time.Second * 5)
		}
		dynamicd, proc, err = dynamic.FindDynamicdProcess(true, time.Second*30)
		if proc != nil {
			util.Info.Println("Running dynamicd process found Pid", proc.Pid)
		} else {
			return fmt.Errorf("Fatal error starting dynamicd daemon %v ", err)
		}
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

		acc, errAccounts := dynamicd.GetMyAccounts(time.Second * 120)
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
		// Create ShutdownApp stuct
		stopWatcher := make(chan struct{})
		dynamic.WatchProcess(stopWatcher, 10, walletpassphrase)
		var sync = false
		closeApp := make(chan struct{})
		stopBridges := make(chan struct{})
		shutdown := rest.AppShutdown{
			Close:       &closeApp,
			StopWatcher: &stopWatcher,
			StopBridges: &stopBridges,
			Dynamicd:    dynamicd,
		}
		// Start Gin web services
		go rest.StartWebServiceRouter(&config, dynamicd, &shutdown, mode)

		al, errLinks := dynamicd.GetActiveLinks(time.Second * 120)
		if errLinks != nil {
			util.Error.Println("GetActiveLinks error", errLinks)
		} else {
			util.Info.Printf("Found %v links\n", len(al.Links))
		}

		go appCommandLoop(&stopBridges, acc, al, &shutdown, dynamicd, status, sync)

		for {
			select {
			case <-*shutdown.Close:
				util.Info.Println("Shutdown close trigger.")
				return nil
			}
		}
	}
	return nil
}
