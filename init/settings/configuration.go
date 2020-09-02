package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

// TODO: Support different ICE service authentication mechanisms

const (
	// ConfigurationFileName is the name of the configuration file
	ConfigurationFileName string = ".webbridge.settings.json"
	// DefaultIceURL is the default ICE service URL
	DefaultIceURL string = "turn:ice.bdap.io:3478"
	// DefaultIceUserName is the default ICE service user name
	DefaultIceUserName string = "test"
	// DefaultIceCredential is the default ICE service credential
	DefaultIceCredential string = "Admin@123"
)

var HomeDir string = ""
var PathSeperator string = ""

// IceServerConfig stores the ICE server configuration information needed for WebRTC connections
type IceServerConfig struct {
	URL        string `json:"URL"`
	UserName   string `json:"UserName"`
	Credential string `json:"Credential"`
}

// Configuration contains the main file settings used when the application launches
type Configuration struct {
	IceServers []IceServerConfig `json:"IceServers"`
}

func isErr(e error) bool {
	if e != nil {
		return true
	}
	return false
}

func (c *Configuration) createDefault() {
	defaultIce := IceServerConfig{
		DefaultIceURL,
		DefaultIceUserName,
		DefaultIceCredential,
	}
	c.IceServers = append(c.IceServers, defaultIce)
	file, _ := json.Marshal(&c)
	err := ioutil.WriteFile(HomeDir+ConfigurationFileName, file, 0644)
	if isErr(err) {
		util.Error.Println("Error writting default configuration file.")
	}
}

// Load reads the configuration file or loads default values
func (c *Configuration) Load(homeDir, pathSeperator string) {
	HomeDir = homeDir
	PathSeperator = pathSeperator
	_, errOpen := os.Open(HomeDir + ConfigurationFileName)
	if isErr(errOpen) {
		util.Error.Println("Configuration file doesn't exist. Creating new configuration with default values.")
		c.createDefault()
	} else {
		file, errRead := ioutil.ReadFile(HomeDir + ConfigurationFileName)
		if isErr(errRead) {
			c.createDefault()
			return
		}
		errUnmarshal := json.Unmarshal([]byte(file), c)
		if isErr(errUnmarshal) {
			util.Error.Println("Error unmarshal configuration file. Overwritting file with default values.")
			c.createDefault()
		}
		util.Info.Println("Configuration loaded.")
	}
}
