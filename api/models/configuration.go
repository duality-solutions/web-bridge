package models

// ConfigurationFile stores the content of the web-bridge configuration file
// swagger:parameters models.ConfigurationFile
type ConfigurationFile struct {
	IceServers   []IceServerConfig `json:"IceServers"`
	WebServer    WebServerConfig
	WalletStatus WalletSetupStatus
}
