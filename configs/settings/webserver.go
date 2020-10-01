package settings

import "github.com/duality-solutions/web-bridge/api/models"

// UpdateWebServer updates an existing web server settings in current running configuration and file
func (c *Configuration) UpdateWebServer(web models.WebServerConfig) bool {
	// check values of web before changing.
	c.mut.Lock()
	defer c.mut.Unlock()
	c.configFile.WebServer = web
	c.updateFile()
	return true
}
