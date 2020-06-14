package dynamic

import (
	"errors"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

func checkBinaryHash(_os, hashDynamicd, hashCli string) error {
	switch _os {
	case "Windows":
		if winDyndHash != hashDynamicd {
			util.Error.Println("Error with", dynamicdName, ". File hash mismatch", winDyndHash, hashDynamicd)
			err := errors.New("Error with " + dynamicdName + ". File hash mismatch " + winDyndHash + " vs " + hashDynamicd)
			return err
		}
		util.Info.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if winDynCliHash != hashCli {
			util.Error.Println("Error with", cliName, ". File hash mismatch", winDynCliHash, hashCli)
			err := errors.New("Error with " + cliName + ". File hash mismatch " + winDynCliHash + " vs " + hashCli)
			return err
		}
		util.Info.Println("File binary hash check pass", cliName, hashCli)
	case "Linux":
		if linuxDyndHash != hashDynamicd {
			util.Error.Println("Error with", dynamicdName, ". File hash mismatch", linuxDyndHash, hashDynamicd)
			err := errors.New("Error with " + dynamicdName + ". File hash mismatch " + linuxDyndHash + " vs " + hashDynamicd)
			return err
		}
		util.Info.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if linuxDynCliHash != hashCli {
			util.Error.Println("Error with", cliName, ". File hash mismatch", linuxDynCliHash, hashCli)
			err := errors.New("Error with " + cliName + ". File hash mismatch " + linuxDynCliHash + " vs " + hashCli)
			return err
		}
		util.Info.Println("File binary hash check pass", cliName, hashCli)
	case "OSX":
		if macDyndHash != hashDynamicd {
			util.Error.Println("Error with", dynamicdName, ". File hash mismatch", macDyndHash, hashDynamicd)
			err := errors.New("Error with " + dynamicdName + ". File hash mismatch " + macDyndHash + " vs " + hashDynamicd)
			return err
		}
		util.Info.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if macDynCliHash != hashCli {
			util.Error.Println("Error with", cliName, ". File hash mismatch", macDynCliHash, hashCli)
			errHash := errors.New("Error with " + cliName + ". File hash mismatch " + macDynCliHash + " vs " + hashCli)
			return errHash
		}
		util.Info.Println("File binary hash check pass", cliName, hashCli)
	}
	return nil
}
