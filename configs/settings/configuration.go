package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/internal/util"
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

var homeDir string = ""
var pathSeperator string = ""

// HomeDir returns the web-bridge home directory
func HomeDir() string {
	return homeDir
}

// PathSeperator returns OS path seperator
func PathSeperator() string {
	return pathSeperator
}

// Configuration contains the main file settings used when the application launches
type Configuration struct {
	mut        *sync.RWMutex
	configFile models.ConfigurationFile
}

func isErr(e error) bool {
	if e != nil {
		return true
	}
	return false
}

func (c *Configuration) updateFile() {
	fileBytes, err := json.Marshal(&c.configFile)
	if err != nil {
		util.Error.Println("Configuration.updateFile() error after marshal configuration file data:", err)
	}
	fileName := (homeDir + ConfigurationFileName)
	err = ioutil.WriteFile(fileName, fileBytes, 0644)
	if err != nil {
		util.Error.Println("Configuration.updateFile() error after WriteFile: ", err)
	}
}

// ToJSON convert the Configuration struct to JSON
func (c *Configuration) ToJSON() models.ConfigurationFile {
	c.mut.Lock()
	defer c.mut.Unlock()
	return c.configFile
}

func (c *Configuration) createDefault() {
	defaultIce := models.IceServerConfig{
		URL:        DefaultIceURL,
		UserName:   DefaultIceUserName,
		Credential: DefaultIceCredential,
	}
	defaultWeb := models.DefaultWebServerConfig()
	defaultWalletStatus := models.DefaultWalletSetupStatus()
	c.configFile.IceServers = append(c.configFile.IceServers, defaultIce)
	c.configFile.WebServer = defaultWeb
	c.configFile.WalletStatus = defaultWalletStatus
	file, _ := json.Marshal(&c)
	err := ioutil.WriteFile(homeDir+ConfigurationFileName, file, 0644)
	if isErr(err) {
		util.Error.Println("Error writting default configuration file.")
	}
}

// Load reads the configuration file or loads default values
func (c *Configuration) Load(dir, seperator string) error {
	c.mut = new(sync.RWMutex)
	homeDir = dir
	pathSeperator = seperator
	_, errOpen := os.Open(homeDir + ConfigurationFileName)
	if isErr(errOpen) {
		util.Error.Println("Configuration file doesn't exist. Creating new configuration with default values.")
		c.createDefault()
	} else {
		file, errRead := ioutil.ReadFile(homeDir + ConfigurationFileName)
		if isErr(errRead) {
			c.createDefault()
			return nil
		}
		errUnmarshal := json.Unmarshal([]byte(file), &c.configFile)
		if isErr(errUnmarshal) {
			util.Error.Println("Error unmarshal configuration file. Overwritting file with default values.")
			c.createDefault()
		}
		if c.configFile.WebServer.AllowCIDR == "" && c.configFile.WebServer.BindAddress == "" && c.configFile.WebServer.ListenPort < 1 {
			c.configFile.WebServer = models.DefaultWebServerConfig()
			c.updateFile()
		} else {
			hasStatus := strings.Index(string(file), `"WalletStatus"`)
			if hasStatus < 0 {
				c.configFile.WalletStatus = models.DefaultWalletSetupStatus()
				c.updateFile()
			}
			if !util.IsValidCIDRList(c.configFile.WebServer.AllowCIDR) {
				return fmt.Errorf("Invalid Web Server allow CIDR: %v", c.configFile.WebServer.AllowCIDR)
			}
			if !util.IsValidIPAddress(c.configFile.WebServer.BindAddress) {
				return fmt.Errorf("Invalid Web Server bind IP address: %v", c.configFile.WebServer.BindAddress)
			}
			if c.configFile.WebServer.ListenPort < 1 {
				return fmt.Errorf("Invalid Web Server port: %v", c.configFile.WebServer.ListenPort)
			}
		}
		util.Info.Println("Configuration loaded successfully.")
	}
	return nil
}
