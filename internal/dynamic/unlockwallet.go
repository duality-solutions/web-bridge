package dynamic

import (
	"fmt"
	"strings"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

// UnlockWallet unlocks the wallet with the given password
func (d *Dynamicd) UnlockWallet(password string) error {
	req, _ := NewRequest("dynamic-cli walletpassphrase \"" + password + "\" 600000")
	response, _ := util.BeautifyJSON(<-d.ExecCmdRequest(req))
	if strings.HasPrefix(response, "null") || strings.Contains(response, "Wallet is already fully unlocked") || strings.Contains(response, "running with an unencrypted wallet") {
		return nil
	}
	return fmt.Errorf("Wallet unlock failed %s", response)
}
