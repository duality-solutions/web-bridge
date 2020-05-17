package webbridge

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

const (
	binaryRepo        string = "https://github.com/duality-solutions/Dynamic"
	binaryReleasePath string = "releases/download"
	binaryVersion     string = "2.4.4.1"
	winDyndHash       string = "NM+nYgDPBk/DTj/BOkG0vdKEEKAHdjRvcNKSeySiUtg="
	winDynCliHash     string = "EoDfzegZT7bFEaHmUr3NMsWYSao0yPpoC6puq1OD8pw="
	linuxDyndHash     string = "8bFnTc9lOMWJsklMFXX4NurK1umSTROLJSvDAmul2MQ="
	linuxDynCliHash   string = "K66Z66XJn+9NEYrUZsqA0UNpGzVHlmEQjlsakioWvn4="
	macDyndHash       string = "AjXMbmI6M1QpKX9JILeMDpdO9d5OkazNKygoRP1y4cg="
	macDynCliHash     string = "4pYr5IQ9NJUrQga7jjhUJ3ThoVQncYGLVv1OyWkRsJs="
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
	ctx            context.Context
	cancelFunc     context.CancelFunc
	datadir        string
	rpcuser        string
	rpcpassword    string
	p2pport        uint16
	rpcport        uint16
	bindaddress    string
	rpcbindaddress string
	cmd            *exec.Cmd
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
) *Dynamicd {
	d := Dynamicd{
		ctx:            ctx,
		cancelFunc:     cancelFunc,
		datadir:        datadir,
		rpcuser:        rpcuser,
		rpcpassword:    rpcpassword,
		p2pport:        p2pport,
		rpcport:        prcport,
		bindaddress:    bindaddress,
		rpcbindaddress: rpcbindaddress,
		cmd:            cmd,
	}
	return &d
}

func loadDynamicd(_os string, archiveExt string) (*Dynamicd, error) {
	if _os == "Windows" {
		dynDir += "\\"
		dynamicdName += ".exe"
		cliName += ".exe"
	} else {
		dynDir += "/"
	}
	if !util.FileExists(dynDir+dynamicdName) || !util.FileExists(dynDir+cliName) {
		fmt.Println(dynamicdName, "or", cliName, "not found. Downloading from Git repo.")
		binPath := dynDir + "dynamic." + archiveExt
		if !util.FileExists(binPath) {
			binaryURL := binaryRepo + "/" + binaryReleasePath + "/v" + binaryVersion + "/Dynamic-" + binaryVersion + "-" + _os + "-" + arch + "." + archiveExt
			fmt.Println("Downloading binaries:", binaryURL)
			err := util.DownloadFile(binPath, binaryURL)
			if err != nil {
				fmt.Println("Binary download failed.", err)
				return nil, err
			}
		} else {
			fmt.Println("Compressed binary found")
		}

		var dir []string
		var errDecompress error
		if _os == "Windows" {
			// unzip archive file
			dir, errDecompress = util.Unzip(binPath, dynDir)
			if errDecompress != nil {
				fmt.Println("Error unzipping file.", binPath, errDecompress)
				return nil, errDecompress
			}
		} else {
			// Extract tar.gz archive file
			dir, errDecompress = util.ExtractTarGz(binPath, dynDir)
			if errDecompress != nil {
				fmt.Println("Error decompressing file.", binPath, errDecompress)
				return nil, errDecompress
			}
		}

		// extract dynamicd.exe dynamid-cli.exe and move
		for _, v := range dir {
			if strings.HasSuffix(v, dynamicdName) {
				fmt.Println("Found", v, "Moving to correct directory")
				errMove := util.MoveFile(v, dynDir+dynamicdName)
				if errMove != nil {
					fmt.Println("Error moving", v, errMove)
					return nil, errMove
				}
				if _os != "Windows" {
					cmd := exec.Command("chmod", "+x", dynDir+dynamicdName)
					errRun := cmd.Run()
					if errRun != nil {
						fmt.Println("Error running chmod for", dynDir+dynamicdName, errRun)
					}
				}
			} else if strings.HasSuffix(v, cliName) {
				fmt.Println("Found", v, "Moving to correct directory")
				errMove := util.MoveFile(v, dynDir+cliName)
				if errMove != nil {
					fmt.Println("Error moving", v, errMove)
					return nil, errMove
				}
				if _os != "Windows" {
					cmd := exec.Command("chmod", "+x", dynDir+cliName)
					errRun := cmd.Run()
					if errRun != nil {
						fmt.Println("Error running chmod for", dynDir+cliName, errRun)
					}
				}
			}
		}
		// clean up extract directory
		dirs, errDirs := util.ListDirectories(dynDir)
		if errDirs != nil {
			fmt.Println("Error listing directories", errDirs)
			return nil, errDirs
		}
		for _, v := range dirs {
			fmt.Println("Deleting directory", dynDir+v)
			errDeleteDir := util.DeleteDirectory(dynDir + v)
			if errDeleteDir != nil {
				fmt.Println("Error deleting directory", v, errDeleteDir)
				return nil, errDeleteDir
			}
		}
		// clean up archive file
		fmt.Println("Cleaning up... Removing unneeded files and directories.")
		if util.FileExists(binPath) {
			fmt.Println("Deleting zip file", binPath)
			errDelete := util.DeleteFile(binPath)
			if errDelete != nil {
				fmt.Println("Error deleting binary archive file", binPath, errDelete)
				return nil, errDelete
			}
		}
	} else {
		fmt.Println("Binaries found", dynamicdName, cliName)
	}
	// check file hashes to make sure they are valid.
	hashDynamicd, _ := util.GetFileHash(dynDir + dynamicdName)
	hashCli, _ := util.GetFileHash(dynDir + cliName)
	var dataDirPath string
	dir, errPath := os.Getwd()
	if errPath != nil {
		return nil, errPath
	}
	switch _os {
	case "Windows":
		if winDyndHash != hashDynamicd {
			fmt.Println("Error with", dynamicdName, ". File hash mismatch", winDyndHash, hashDynamicd)
			errHash := errors.New("Error with " + dynamicdName + ". File hash mismatch " + winDyndHash + " vs " + hashDynamicd)
			return nil, errHash
		}
		fmt.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if winDynCliHash != hashCli {
			fmt.Println("Error with", cliName, ". File hash mismatch", winDynCliHash, hashCli)
			errHash := errors.New("Error with " + cliName + ". File hash mismatch " + winDynCliHash + " vs " + hashCli)
			return nil, errHash
		}
		fmt.Println("File binary hash check pass", cliName, hashCli)
		dataDirPath = dir + "\\" + dynDir + ".dynamic"
	case "Linux":
		if linuxDyndHash != hashDynamicd {
			fmt.Println("Error with", dynamicdName, ". File hash mismatch", linuxDyndHash, hashDynamicd)
			errHash := errors.New("Error with " + dynamicdName + ". File hash mismatch " + linuxDyndHash + " vs " + hashDynamicd)
			return nil, errHash
		}
		fmt.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if linuxDynCliHash != hashCli {
			fmt.Println("Error with", cliName, ". File hash mismatch", linuxDynCliHash, hashCli)
			errHash := errors.New("Error with " + cliName + ". File hash mismatch " + linuxDynCliHash + " vs " + hashCli)
			return nil, errHash
		}
		fmt.Println("File binary hash check pass", cliName, hashCli)
		dataDirPath = dir + "/" + dynDir + ".dynamic"
	case "OSX":
		if macDyndHash != hashDynamicd {
			fmt.Println("Error with", dynamicdName, ". File hash mismatch", macDyndHash, hashDynamicd)
			errHash := errors.New("Error with " + dynamicdName + ". File hash mismatch " + macDyndHash + " vs " + hashDynamicd)
			return nil, errHash
		}
		fmt.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if macDynCliHash != hashCli {
			fmt.Println("Error with", cliName, ". File hash mismatch", macDynCliHash, hashCli)
			errHash := errors.New("Error with " + cliName + ". File hash mismatch " + macDynCliHash + " vs " + hashCli)
			return nil, errHash
		}
		fmt.Println("File binary hash check pass", cliName, hashCli)
		dataDirPath = dir + "/" + dynDir + ".dynamic"
	}
	if !util.DirectoryExists(dataDirPath) {
		errMkdir := os.Mkdir(dataDirPath, 0755)
		if errMkdir != nil {
			return nil, errMkdir
		}
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
	dynamicd := newDynamicd(ctx, cancel, dataDirPath, rpcUser, rpcPassword, defaultPort, defaultRPCPort, defaultBind, defaultBind, cmd)

	return dynamicd, nil
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
