package dynamic

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/duality-solutions/web-bridge/blockchain/rpc"
	"github.com/duality-solutions/web-bridge/configs/settings"
	"github.com/duality-solutions/web-bridge/internal/util"
)

const (
	binaryRepo        string = "https://github.com/duality-solutions/Dynamic"
	binaryReleasePath string = "releases/download"
	binaryVersion     string = "2.4.4.2"
	winDyndHash       string = "KuHhwzO4OinWd6AD5xRAoDYG+T7RhNISXAocSBd+XlQ="
	winDynCliHash     string = "AdSfB1vGR9dOsuAWNyQ0yJ3m5rsypKrF8lMmdd5/qec="
	linuxDyndHash     string = "p/ruW+P3hGeq6+lc7l9B1bRQRbEIm3WelfM9I7ws4Io="
	linuxDynCliHash   string = "cMN/8n5Pguq0j305Lm7iEHohfIgrhIxWpvlFkGf00Js="
	macDyndHash       string = ""
	macDynCliHash     string = ""

	defaultName    string = DefaultName
	defaultPort    uint16 = DefaultPort
	defaultRPCPort uint16 = DefaultRPCPort
	defaultBind    string = "127.0.0.1"
	arch           string = "x64"
)

var initVars bool = false
var dynDir string = "dynamic"
var dynamicdName string = defaultName
var cliName string = "dynamic-cli"

// Dynamicd stores information about the binded dynamicd process
type Dynamicd struct {
	Ctx            context.Context
	CancelFunc     context.CancelFunc
	Datadir        string
	RPCUser        string
	RPCPassword    string
	P2PPort        uint16
	RPCPort        uint16
	BindAddress    string
	RPCBindAddress string
	Cmd            *exec.Cmd
	ConfigRPC      rpc.Config
	WalletPassword string
}

func newDynamicd(
	ctx context.Context,
	cancelFunc context.CancelFunc,
	datadir string,
	rpcuser string,
	rpcpassword string,
	p2pport uint16,
	prcport uint16,
	bindaddress string,
	rpcbindaddress string,
	cmd *exec.Cmd,
	configRPC rpc.Config,
) *Dynamicd {
	d := Dynamicd{
		Ctx:            ctx,
		CancelFunc:     cancelFunc,
		Datadir:        datadir,
		RPCUser:        rpcuser,
		RPCPassword:    rpcpassword,
		P2PPort:        p2pport,
		RPCPort:        prcport,
		BindAddress:    bindaddress,
		RPCBindAddress: rpcbindaddress,
		Cmd:            cmd,
		ConfigRPC:      configRPC,
	}
	return &d
}

func getCmd(ctx context.Context, dataDirPath, binpath, rpcUser, rpcPassword string) *exec.Cmd {
	return exec.CommandContext(ctx, binpath+dynamicdName,
		"-datadir="+dataDirPath,
		"-port="+strconv.Itoa(int(defaultPort)),
		"-rpcuser="+rpcUser,
		"-rpcpassword="+rpcPassword,
		"-rpcport="+strconv.Itoa(int(defaultRPCPort)),
		"-rpcbind="+defaultBind,
		"-rpcallowip="+defaultBind,
		"-server=1",
		"-stealthtx=1",
		"-txindex=1",
		"-addressindex=1",
		"-timestampindex=1",
		"-spentindex=1",
		"-addnode=206.189.68.224",
		"-addnode=138.197.193.115",
		"-addnode=dynexplorer.duality.solutions")
}

func loadDynamicd(_os, archiveExt string) (*Dynamicd, error) {
	var dataDirPath string = settings.HomeDir()
	if !util.DirectoryExists(dataDirPath) {
		err := util.AddDirectory(dataDirPath)
		if err != nil {
			return nil, fmt.Errorf("loadDynamicd failed after AddDirectory root %v. %v", dataDirPath, err)
		}
	}
	var dirDelimit string = settings.PathSeperator()
	if _os == "Windows" {
		if !initVars {
			dynDir += dirDelimit
			dynamicdName += ".exe"
			cliName += ".exe"
			initVars = true
		}

	} else {
		if !initVars {
			dynDir += dirDelimit
			initVars = true
		}
	}
	err := downloadBinaries(_os, dataDirPath, dynamicdName, cliName, archiveExt)
	if err != nil {
		return nil, fmt.Errorf("loadDynamicd failed at downloadBinaries. %v", err)
	}
	// check file hashes to make sure they are valid.
	hashDynamicd, _ := util.GetFileHash(dataDirPath + dynamicdName)
	hashCli, _ := util.GetFileHash(dataDirPath + cliName)
	err = checkBinaryHash(_os, hashDynamicd, hashCli)
	if err != nil {
		return nil, fmt.Errorf("loadDynamicd failed at checkBinaryHash. %v", err)
	}
	// check to make sure .dynamic directory exists and create if doesn't
	dir := dataDirPath
	rpcUser, _ := util.RandomString(12)
	rpcPassword, _ := util.RandomString(18)
	dataDirPath = dir + "." + dynDir
	ctx, cancel := context.WithCancel(context.Background())
	var cmd *exec.Cmd
	if !util.DirectoryExists(dataDirPath) {
		err = util.AddDirectory(dataDirPath)
		if err != nil {
			defer cancel()
			return nil, fmt.Errorf("loadDynamicd failed after AddDirectory %v. %v", dataDirPath, err)
		}
		cmd = getCmd(ctx, dataDirPath, dir, rpcUser, rpcPassword)
	} else {
		// read username and password from config file
		var userFound, passFound bool = false, false
		configPath := dataDirPath + dirDelimit + "dynamic.conf"
		conf, err := GetDynamicConfig(configPath)
		if err == nil {
			user, err := ParseDynamicConfigValue(conf, "rpcuser=")
			if err == nil {
				rpcUser = user
				userFound = true
			} else {
				util.Error.Println("loadDynamicd error after ParseDynamicConfValue rpcUser", err)
			}
			pass, err := ParseDynamicConfigValue(conf, "rpcpassword=")
			if err == nil {
				rpcPassword = pass
				passFound = true
			} else {
				util.Error.Println("loadDynamicd error after ParseDynamicConfValue rpcPassword", err)
			}
		} else {
			util.Error.Println("loadDynamicd error after GetDynamicConf", err)
		}
		if userFound && passFound {
			cmd = getCmd(ctx, dataDirPath, dir, rpcUser, rpcPassword)
		} else {
			cmd = getCmd(ctx, dataDirPath, dir, rpcUser, rpcPassword)
		}
	}

	configRPC := rpc.Config{
		RPCServer:   defaultBind + ":" + strconv.Itoa(int(defaultRPCPort)),
		RPCUser:     rpcUser,
		RPCPassword: rpcPassword,
		NoTLS:       true,
	}
	util.Info.Println("dynamicd starting...")
	time.Sleep(time.Second * 5)
	dynamicd := newDynamicd(ctx, cancel, dataDirPath, rpcUser, rpcPassword, defaultPort, defaultRPCPort, defaultBind, defaultBind, cmd, configRPC)
	if errStart := dynamicd.Cmd.Start(); errStart != nil {
		return nil, fmt.Errorf("loadDynamicd failed at dynamicd.Cmd.Start(). %v", errStart)
	}
	time.Sleep(time.Second * 5)
	util.Info.Println("dynamicd started")
	return dynamicd, nil
}

// ExecCmdRequest runs a dynamic-cli command using the RPCRequest struct
func (d *Dynamicd) ExecCmdRequest(req *RPCRequest) <-chan string {
	c := make(chan string)
	go func() {
		byteCmd, _ := json.Marshal(req)
		byteResp, errResp := rpc.SendPostRequest(byteCmd, &d.ConfigRPC)
		if errResp != nil {
			c <- errResp.Error()
		} else {
			c <- string(byteResp)
		}
	}()
	return c
}

// LoadRPCDynamicd is used to create and run a managed dynamicd full node and cli.
func LoadRPCDynamicd() (*Dynamicd, error) {
	var dynamicd *Dynamicd
	var err error
	switch runtime.GOOS {
	case "windows":
		dynamicd, err = loadDynamicd("Windows", "zip")
		if err != nil {
			return nil, err
		}
	case "linux":
		dynamicd, err = loadDynamicd("Linux", "tar.gz")
		if err != nil {
			return nil, err
		}
	case "darwin":
		dynamicd, err = loadDynamicd("OSX", "tar.gz")
		if err != nil {
			return nil, err
		}
	}
	return dynamicd, nil
}
