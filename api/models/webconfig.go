package models

import "strconv"

// WebServerConfig is used to store the built in Web Server configurations
// swagger:parameters models.WebServerConfig
type WebServerConfig struct {
	BindAddress string
	ListenPort  uint16
	AllowCIDR   string
}

// WebServerRestartRequest tells the server when it should restart
// swagger:parameters models.WebServerRestartRequest
type WebServerRestartRequest struct {
	RestartEpoch int64 `json:"restart_epoch"`
}

func DefaultWebServerConfig() WebServerConfig {
	return WebServerConfig{
		BindAddress: "127.0.0.1",
		ListenPort:  35350,
		AllowCIDR:   "127.0.0.0/8, ::1/128",
	}
}

func (w WebServerConfig) PortString() string {
	return strconv.Itoa(int(w.ListenPort))
}

func (w WebServerConfig) AddressPortString() string {
	address := w.BindAddress
	if address == "127.0.0.1" {
		address = "localhost"
	}
	return address + ":" + w.PortString()
}

func (w WebServerConfig) AddressPortRawString() string {
	return w.BindAddress + ":" + w.PortString()
}
