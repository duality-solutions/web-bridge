package models

import "strconv"

// WebServerConfig is used to store the built in Web Server configurations
// swagger:parameters models.WebServerConfig
type WebServerConfig struct {
	BindAddress string
	ListenPort  uint16
	AllowCIDR   string
}

func DefaultWebServerConfig() WebServerConfig {
	return WebServerConfig{
		BindAddress: "127.0.0.1",
		ListenPort:  35350,
		AllowCIDR:   "127.0.0.1/0",
	}
}

func (w WebServerConfig) PortString() string {
	return strconv.Itoa(int(w.ListenPort))
}

func (w WebServerConfig) AddressPortString() string {
	return w.BindAddress + ":" + w.PortString()
}
