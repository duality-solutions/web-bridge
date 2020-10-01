package settings

import "github.com/duality-solutions/web-bridge/api/models"

// IceServers returns the current configuration ICE servers
func (c *Configuration) IceServers() *[]models.IceServerConfig {
	c.mut.Lock()
	defer c.mut.Unlock()
	return &c.configFile.IceServers
}

// AddIceServer adds a new ICE Server to the current running configuration and file
func (c *Configuration) AddIceServer(ice models.IceServerConfig) bool {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.configFile.IceServers = append(c.configFile.IceServers, ice)
	c.updateFile()
	return true
}

// UpdateIceServer updates an existing ICE Server in current running configuration and file
func (c *Configuration) UpdateIceServer(index int, ice models.IceServerConfig) bool {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.configFile.IceServers[index] = ice
	c.updateFile()
	return true
}

// DeleteIceServer deleted an existing ICE Server from current running configuration and file
func (c *Configuration) DeleteIceServer(index int) bool {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.configFile.IceServers = append(c.configFile.IceServers[:index], c.configFile.IceServers[index+1:]...)
	c.updateFile()
	return true
}

// ReplaceIceServers deleted an existing ICE Server from current running configuration and file
func (c *Configuration) ReplaceIceServers(fileData models.ConfigurationFile) bool {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.configFile.IceServers = fileData.IceServers
	c.updateFile()
	return true
}
