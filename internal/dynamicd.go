package webbridge

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	util "github.com/amirabrams/webbridge/internal/utilities"
)

const (
	binaryRepo        string = "https://github.com/duality-solutions/Dynamic"
	binaryVersionPath string = "releases/download/v2.4.4.1"
	winDyndHash       string = "NM+nYgDPBk/DTj/BOkG0vdKEEKAHdjRvcNKSeySiUtg="
	winDynCliHash     string = "EoDfzegZT7bFEaHmUr3NMsWYSao0yPpoC6puq1OD8pw="
	linuxDyndHash     string = ""
	linuxDynCliHash   string = ""
	macDyndHash       string = ""
	macDynCliHash     string = ""
)

/*
docker volume create --name=dynamicd-data
docker run -e DISABLEWALLET=0 -v dynamicd-data:/dynamic --name=dynamicd -d -p 33300:33300 -p 127.0.0.1:33350:33350 dualitysolutions/docker-dynamicd
*/

// LoadDockerDynamicd is used to create and run a Docker dynamicd full node and cli.
func LoadDockerDynamicd() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "dynamicd",
		ExposedPorts: nat.PortSet{"33350": struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("33350"): {{HostIP: "127.0.0.1", HostPort: "8080"}}},
	}, nil, "dynamicd")
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	return cli, nil
}

func loadWindowsDynamicd() {
	//https://github.com/duality-solutions/Dynamic/releases/tag/v2.4.4.1/Dynamic-2.4.4.1-Windows-x64.zip
	// look for binary.
	// if not found, download a new one.
	dynDir := "dynamic\\"
	dynamicdName := "dynamicd.exe"
	cliName := "dynamic-cli.exe"
	if !util.FileExists("dynamic\\dynamicd.exe") || !util.FileExists("dynamic\\dynamic-cli.exe") {
		fmt.Println("dynamicd.exe or dynamid-cli.exe not found. Downloading from repo.")
		binZipPath := "dynamic\\dynamic-bin.zip"
		if !util.FileExists(binZipPath) {
			binaryURL := binaryRepo + "/" + binaryVersionPath + "/Dynamic-2.4.4.1-Windows-x64.zip"
			fmt.Println("Downloading binaries:", binaryURL)
			err := util.DownloadFile(binZipPath, binaryURL)
			if err != nil {
				fmt.Println("Binary download failed.", err)
			}
		} else {
			fmt.Println("Compressed binary found")
		}
		// unzip
		dir, err := util.Unzip(binZipPath, dynDir)
		if err != nil {
			fmt.Println("Error unzipping file.", binZipPath, err)
		}
		// extract dynamicd.exe dynamid-cli.exe and move
		for _, v := range dir {
			if strings.HasSuffix(v, dynamicdName) {
				fmt.Println("Found", v, "Moving to correct directory")
				errMove := util.MoveFile(v, dynDir+dynamicdName)
				if errMove != nil {
					fmt.Println("Error moving", v, errMove)
				}
			} else if strings.HasSuffix(v, cliName) {
				fmt.Println("Found", v, "Moving to correct directory")
				errMove := util.MoveFile(v, dynDir+cliName)
				if errMove != nil {
					fmt.Println("Error moving", v, errMove)
				}
			}
		}
		// clean up
		fmt.Println("Cleaning up... Removing unneeded files and directories.")
		if util.FileExists(binZipPath) {
			fmt.Println("Deleting zip file", binZipPath)
			errDelete := util.DeleteFile(binZipPath)
			if errDelete != nil {
				fmt.Println("Error deleting binary archive file", binZipPath, errDelete)
			}
		}
		// clean up extract directory
		dirs, errDirs := util.ListDirectories(dynDir)
		if errDirs != nil {
			fmt.Println("Error listing directories", errDirs)
		}
		for _, v := range dirs {
			fmt.Println("Deleting directory", dynDir+v)
			errDeleteDir := util.DeleteDirectory(dynDir + v)
			if errDeleteDir != nil {
				fmt.Println("Error deleting directory", v, errDeleteDir)
			}
		}
	}
	// check file hashes to make sure they are valid.
	hashDynamicd, _ := util.GetFileHash(dynDir + dynamicdName)
	hashCli, _ := util.GetFileHash(dynDir + cliName)

	if winDyndHash != hashDynamicd {
		fmt.Println("Error with", dynamicdName, ". File hash mismatch", winDyndHash, hashDynamicd)
		// TODO panic
	} else {
		fmt.Println("File binary hash check pass", dynamicdName, hashDynamicd)
	}
	if winDynCliHash != hashCli {
		fmt.Println("Error with", cliName, ". File hash mismatch", winDynCliHash, hashCli)
		// TODO panic
	} else {
		fmt.Println("File binary hash check pass", cliName, hashCli)
	}
}

func loadLinuxDynamicd() {
	//https://github.com/duality-solutions/Dynamic/releases/tag/v2.4.4.1/Dynamic-2.4.4.1-Linux-x64.tar.gz

}

func loadMacOSDynamicd() {
	//https://github.com/duality-solutions/Dynamic/releases/tag/v2.4.4.1/Dynamic-2.4.4.1-OSX-x64.tar.gz

}

// LoadRPCDynamicd is used to create and run a managed dynamicd full node and cli.
func LoadRPCDynamicd() error {
	switch runtime.GOOS {
	case "windows":
		loadWindowsDynamicd()
	case "linux":
		loadLinuxDynamicd()
	case "darwin":
		loadMacOSDynamicd()
	}
	return nil
}
