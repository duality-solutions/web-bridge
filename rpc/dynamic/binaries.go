package dynamic

import (
	"os/exec"
	"strings"

	"github.com/duality-solutions/web-bridge/internal/util"
)

func downloadBinaries(_os, dynDir, dynamicName, cliName, archiveExt string) error {
	if !util.FileExists(dynDir+dynamicdName) || !util.FileExists(dynDir+cliName) {
		util.Info.Println(dynamicdName, "or", cliName, "not found. Downloading from Git repo.")
		binPath := dynDir + "dynamic." + archiveExt
		if !util.FileExists(binPath) {
			binaryURL := binaryRepo + "/" + binaryReleasePath + "/v" + binaryVersion + "/Dynamic-" + binaryVersion + "-" + _os + "-" + arch + "." + archiveExt
			util.Info.Println("Downloading binaries:", binaryURL)
			err := util.DownloadFile(binPath, binaryURL)
			if err != nil {
				util.Error.Println("Binary download failed.", err)
				return err
			}
		} else {
			util.Info.Println("Compressed binary found")
		}

		var dir []string
		var errDecompress error
		if _os == "Windows" {
			// unzip archive file
			dir, errDecompress = util.Unzip(binPath, dynDir)
			if errDecompress != nil {
				util.Error.Println("Error unzipping file.", binPath, errDecompress)
				return errDecompress
			}
		} else {
			// Extract tar.gz archive file
			dir, errDecompress = util.ExtractTarGz(binPath, dynDir)
			if errDecompress != nil {
				util.Error.Println("Error decompressing file.", binPath, errDecompress)
				return errDecompress
			}
		}

		// extract dynamicd.exe dynamid-cli.exe and move
		for _, v := range dir {
			if strings.HasSuffix(v, dynamicdName) {
				util.Info.Println("Found", v, "Moving to correct directory")
				errMove := util.MoveFile(v, dynDir+dynamicdName)
				if errMove != nil {
					util.Error.Println("Error moving", v, errMove)
					return errMove
				}
				if _os != "Windows" {
					cmd := exec.Command("chmod", "+x", dynDir+dynamicdName)
					errRun := cmd.Run()
					if errRun != nil {
						util.Error.Println("Error running chmod for", dynDir+dynamicdName, errRun)
					}
				}
			} else if strings.HasSuffix(v, cliName) {
				util.Info.Println("Found", v, "Moving to correct directory")
				errMove := util.MoveFile(v, dynDir+cliName)
				if errMove != nil {
					util.Error.Println("Error moving", v, errMove)
					return errMove
				}
				if _os != "Windows" {
					cmd := exec.Command("chmod", "+x", dynDir+cliName)
					errRun := cmd.Run()
					if errRun != nil {
						util.Error.Println("Error running chmod for", dynDir+cliName, errRun)
					}
				}
			}
		}
		// clean up extract directory
		dirs, errDirs := util.ListDirectories(dynDir)
		if errDirs != nil {
			util.Error.Println("Error listing directories", errDirs)
			return errDirs
		}
		for _, v := range dirs {
			if !strings.HasPrefix(v, ".dynamic") {
				util.Info.Println("Deleting directory", dynDir+v)
				errDeleteDir := util.DeleteDirectory(dynDir + v)
				if errDeleteDir != nil {
					util.Error.Println("Error deleting directory", v, errDeleteDir)
					return errDeleteDir
				}
			}
		}
		// clean up archive file
		util.Info.Println("Cleaning up... Removing unneeded files and directories.")
		if util.FileExists(binPath) {
			util.Info.Println("Deleting zip file", binPath)
			errDelete := util.DeleteFile(binPath)
			if errDelete != nil {
				util.Error.Println("Error deleting binary archive file", binPath, errDelete)
				return errDelete
			}
		}
	} else {
		util.Info.Println("Binaries found", dynamicdName, cliName)
	}
	return nil
}
