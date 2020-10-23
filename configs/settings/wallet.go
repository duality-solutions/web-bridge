package settings

import "github.com/duality-solutions/web-bridge/api/models"

// WalletSetupStatus returns the wallet setup status in the running configuration
func (c *Configuration) WalletSetupStatus() *models.WalletSetupStatus {
	c.mut.Lock()
	defer c.mut.Unlock()
	return &c.configFile.WalletStatus
}
