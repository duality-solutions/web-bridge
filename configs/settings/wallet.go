package settings

import "github.com/duality-solutions/web-bridge/api/models"

// WalletSetupStatus returns the wallet setup status in the running configuration
func (c *Configuration) WalletSetupStatus() *models.WalletSetupStatus {
	c.mut.Lock()
	defer c.mut.Unlock()
	return &c.configFile.WalletStatus
}

// UpdateWalletSetup updates the current running wallet setup status configration
func (c *Configuration) UpdateWalletSetup(setup models.WalletSetupStatus) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.configFile.WalletStatus = setup
	UnlockedUntil := c.configFile.WalletStatus.UnlockedUntil
	// we don't want to save the unlock until epoch time in the configuration file.
	c.configFile.WalletStatus.UnlockedUntil = 0
	c.updateFile()
	// set the running value back to the original value
	c.configFile.WalletStatus.UnlockedUntil = UnlockedUntil
	return
}
