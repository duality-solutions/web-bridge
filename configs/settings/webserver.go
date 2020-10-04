package settings

import (
	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/internal/util"
)

// WebServer returns the existing web server settings in current running configuration and file
func (c *Configuration) WebServer() models.WebServerConfig {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.configFile.WebServer
}

// UpdateWebServer updates an existing web server settings in current running configuration and file
func (c *Configuration) UpdateWebServer(web models.WebServerConfig) bool {
	// check values of web before changing.
	if web.BindAddress == "localhost" {
		web.BindAddress = "127.0.0.1"
	}
	if !util.IsValidIPAddress(web.BindAddress) {
		return false
	}
	if !util.IsValidCIDRList(web.AllowCIDR) {
		return false
	}
	if web.ListenPort < 1 {
		return false
	}
	c.mut.Lock()
	defer c.mut.Unlock()
	c.configFile.WebServer = web
	c.updateFile()
	return true
}
