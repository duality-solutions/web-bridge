package dynamic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	rpcclient "github.com/duality-solutions/web-bridge/rpc"
)

const (
	binaryRepo        string = "https://github.com/duality-solutions/Dynamic"
	binaryReleasePath string = "releases/download"
	binaryVersion     string = "2.4.4.2"
	winDyndHash       string = "DNIgNNmzIi3BG1vLHzzkwnYTWy8Tyt0L6/dh//p6eNg="
	winDynCliHash     string = "JELpv7Mz6+axrRkMfDgRtEknvuAWiSLa4fPz7ahN6TU="
	linuxDyndHash     string = "K96nSCIM40zcBFhrYCdro9wwhXQFvdaY/KxHx0i7h/k="
	linuxDynCliHash   string = "t5FOoZJoZ8nogTh5Qfyr62/PWXHXqdFmufYovuzBQWU="
	macDyndHash       string = ""
	macDynCliHash     string = ""
)

var dynamicdName string = "dynamicd"
var cliName string = "dynamic-cli"
var dynDir string = "dynamic"
var arch = "x64"
var defaultPort uint16 = 33334
var defaultRPCPort uint16 = 33335
var defaultBind string = "127.0.0.1"

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
	ConfigRPC      rpcclient.Config
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
	configRPC rpcclient.Config,
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

func loadDynamicd(_os string, archiveExt string) (*Dynamicd, error) {
	var dataDirPath string
	var dirDelimit string
	if _os == "Windows" {
		dirDelimit = "\\"
		dynDir += dirDelimit
		dynamicdName += ".exe"
		cliName += ".exe"
	} else {
		dirDelimit = "/"
		dynDir += dirDelimit
	}
	err := downloadBinaries(_os, dynDir, dynamicdName, cliName, archiveExt)
	if err != nil {
		return nil, err
	}
	// check file hashes to make sure they are valid.
	hashDynamicd, _ := util.GetFileHash(dynDir + dynamicdName)
	hashCli, _ := util.GetFileHash(dynDir + cliName)
	err = checkBinaryHash(_os, hashDynamicd, hashCli)
	if err != nil {
		return nil, err
	}
	// check to make sure .dynamic directory exists and create if doesn't
	dir, errPath := os.Getwd()
	if errPath != nil {
		return nil, errPath
	}
	dataDirPath = dir + dirDelimit + dynDir + ".dynamic"
	err = util.AddDirectory(dataDirPath)
	if err != nil {
		return nil, err
	}

	rpcUser, errUser := util.RandomString(12)
	if errUser != nil {
		return nil, errUser
	}
	rpcPassword, errPassword := util.RandomString(18)
	if errPassword != nil {
		return nil, errPassword
	}
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, dynDir+dynamicdName,
		"-datadir="+dataDirPath,
		"-port="+string(defaultPort),
		"-rpcuser="+rpcUser,
		"-rpcpassword="+rpcPassword,
		"-rpcport="+strconv.Itoa(int(defaultRPCPort)),
		"-rpcbind="+defaultBind,
		"-rpcallowip="+defaultBind,
		"-server=1",
		"-addnode=206.189.68.224",
		"-addnode=138.197.193.115",
		"-addnode=dynexplorer.duality.solutions",
	)
	configRPC := rpcclient.Config{
		RPCServer:   defaultBind + ":" + strconv.Itoa(int(defaultRPCPort)),
		RPCUser:     rpcUser,
		RPCPassword: rpcPassword,
		NoTLS:       true,
	}
	fmt.Println("dynamicd starting...")
	dynamicd := newDynamicd(ctx, cancel, dataDirPath, rpcUser, rpcPassword, defaultPort, defaultRPCPort, defaultBind, defaultBind, cmd, configRPC)
	if errStart := dynamicd.Cmd.Start(); errStart != nil {
		return nil, errStart
	}
	time.Sleep(time.Second * 5)
	fmt.Println("dynamicd started")
	return dynamicd, nil
}

// ExecCmdRequest runs a dynamic-cli command using the RPCRequest struct
func (d *Dynamicd) ExecCmdRequest(req *RPCRequest) <-chan string {
	c := make(chan string)
	go func() {
		byteCmd, _ := json.Marshal(req)
		byteResp, errResp := rpcclient.SendPostRequest(byteCmd, &d.ConfigRPC)
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
