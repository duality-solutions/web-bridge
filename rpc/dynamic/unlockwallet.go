package dynamic

import (
	"fmt"
	"strings"
	"time"
)

func checkPrefix(r string) bool {
	if strings.HasPrefix(r, "null") || strings.Contains(r, "Wallet is already fully unlocked") || strings.Contains(r, "running with an unencrypted wallet") {
		return true
	}
	return false
}

func passwordSuccessful(r string) bool {
	if strings.HasPrefix(r, "null") {
		return true
	}
	return false
}

// UnlockWallet unlocks the wallet with the given password
func (d *Dynamicd) UnlockWallet(password string) error {
	req, _ := NewRequest("dynamic-cli walletpassphrase \"" + password + "\" 600000")
	response, _ := <-d.ExecCmdRequest(req)
	if checkPrefix(response) {
		if passwordSuccessful(response) {
			d.WalletPassword = password
		}
		return nil
	}
	var loadingMessage string = "Loading wallet..."
	if strings.Contains(response, loadingMessage) {
		fmt.Println(loadingMessage)
		time.Sleep(time.Second * 5)
		for strings.Contains(response, loadingMessage) {
			response, _ = <-d.ExecCmdRequest(req)
			if checkPrefix(response) {
				if passwordSuccessful(response) {
					d.WalletPassword = password
				}
				return nil
			}
			fmt.Println(loadingMessage)
			time.Sleep(time.Second * 5)
		}
	}
	return fmt.Errorf("Wallet unlock failed %s", response)
}
