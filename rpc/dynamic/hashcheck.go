package dynamic

import (
	"errors"
	"fmt"
)

func checkBinaryHash(_os, hashDynamicd, hashCli string) error {
	switch _os {
	case "Windows":
		if winDyndHash != hashDynamicd {
			fmt.Println("Error with", dynamicdName, ". File hash mismatch", winDyndHash, hashDynamicd)
			err := errors.New("Error with " + dynamicdName + ". File hash mismatch " + winDyndHash + " vs " + hashDynamicd)
			return err
		}
		fmt.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if winDynCliHash != hashCli {
			fmt.Println("Error with", cliName, ". File hash mismatch", winDynCliHash, hashCli)
			err := errors.New("Error with " + cliName + ". File hash mismatch " + winDynCliHash + " vs " + hashCli)
			return err
		}
		fmt.Println("File binary hash check pass", cliName, hashCli)
	case "Linux":
		if linuxDyndHash != hashDynamicd {
			fmt.Println("Error with", dynamicdName, ". File hash mismatch", linuxDyndHash, hashDynamicd)
			err := errors.New("Error with " + dynamicdName + ". File hash mismatch " + linuxDyndHash + " vs " + hashDynamicd)
			return err
		}
		fmt.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if linuxDynCliHash != hashCli {
			fmt.Println("Error with", cliName, ". File hash mismatch", linuxDynCliHash, hashCli)
			err := errors.New("Error with " + cliName + ". File hash mismatch " + linuxDynCliHash + " vs " + hashCli)
			return err
		}
		fmt.Println("File binary hash check pass", cliName, hashCli)
	case "OSX":
		if macDyndHash != hashDynamicd {
			fmt.Println("Error with", dynamicdName, ". File hash mismatch", macDyndHash, hashDynamicd)
			err := errors.New("Error with " + dynamicdName + ". File hash mismatch " + macDyndHash + " vs " + hashDynamicd)
			return err
		}
		fmt.Println("File binary hash check pass", dynamicdName, hashDynamicd)
		if macDynCliHash != hashCli {
			fmt.Println("Error with", cliName, ". File hash mismatch", macDynCliHash, hashCli)
			errHash := errors.New("Error with " + cliName + ". File hash mismatch " + macDynCliHash + " vs " + hashCli)
			return errHash
		}
		fmt.Println("File binary hash check pass", cliName, hashCli)
	}
	return nil
}
