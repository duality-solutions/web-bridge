package models

// WalletSetupStatus stores status of the wallet setup
// swagger:parameters models.ConfigurationFile
type WalletSetupStatus struct {
	MnemonicBackup  bool  `json:"MnemonicBackup"`
	HasAccounts     bool  `json:"HasAccounts"`
	HasLinks        bool  `json:"HasLinks"`
	HasTransactions bool  `json:"HasTransactions"`
	WalletEncrypted bool  `json:"WalletEncrypted"`
	UnlockedUntil   int64 `json:"UnlockedUntil"`
}

// DefaultWalletSetupStatus creates a default WalletSetupStatus struct
func DefaultWalletSetupStatus() WalletSetupStatus {
	return WalletSetupStatus{
		MnemonicBackup:  false,
		HasAccounts:     false,
		HasLinks:        false,
		HasTransactions: false,
		WalletEncrypted: false,
	}
}
