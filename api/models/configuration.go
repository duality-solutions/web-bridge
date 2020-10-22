package models

// WalletSetupStatus stores status of the wallet setup
// swagger:parameters models.ConfigurationFile
type WalletSetupStatus struct {
	MnemonicBackup  bool `json:"MnemonicBackup"`
	HasAccounts     bool `json:"HasAccounts"`
	HasLinks        bool `json:"HasLinks"`
	HasTransactions bool `json:"HasTransactions"`
}

// DefaultWalletSetupStatus creates a default WalletSetupStatus struct
func DefaultWalletSetupStatus() WalletSetupStatus {
	return WalletSetupStatus{
		MnemonicBackup:  false,
		HasAccounts:     false,
		HasLinks:        false,
		HasTransactions: false,
	}
}

// ConfigurationFile stores the content of the web-bridge configuration file
// swagger:parameters models.ConfigurationFile
type ConfigurationFile struct {
	IceServers   []IceServerConfig `json:"IceServers"`
	WebServer    WebServerConfig
	WalletStatus WalletSetupStatus
}
