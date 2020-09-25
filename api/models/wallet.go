package models

type UnlockWalletRequest struct {
	Passphrase string `json:"passphrase"`
	Timeout    int    `json:"timeout"`
	MixingOnly bool   `json:"mixingonly"`
}

type EncryptWalletRequest struct {
	Passphrase string `json:"passphrase"`
}

type ChangePassphraseRequest struct {
	OldPassphrase string `json:"oldpassphrase"`
	NewPassphrase string `json:"newpassphrase"`
}

// MnemonicResponse contains the wallet HD seed and mnemonic information
// swagger:parameters models.WalletSeedResponse
type MnemonicResponse struct {
	// HDSeed (string, required) deterministic wallet seed
	HDSeed string `json:"hdseed"`
	// Mnemonic (string, required) mnemonic associated with HD seed
	Mnemonic string `json:"mnemonic"`
	// MnemonicPassphrase (string, optional)  mnemonic passphrase used as the 13th or 25th word
	MnemonicPassphrase string `json:"mnemonicpassphrase"`
}

// ImportMnemonicRequest request payload used to import mnemonic
// swagger:parameters models.ImportMnemonicRequest
type ImportMnemonicRequest struct {
	// Mnemonic (string, required) mnemonic delimited by the dash charactor (-) or space. Use 12 or 24 words
	Mnemonic string `json:"mnemonic"`
	// Language (string, optional: english|french|chinesesimplified|chinesetraditional|italian|japanese|korean|spanish)
	Language string `json:"language"`
	// Passphrase (string, optional) mnemonic passphrase used as the 13th or 25th word
	Passphrase string `json:"passphrase"`
}

/*
type HDAccount struct {
	HdAccountIndex     int `json:"hdaccountindex"`
	HdExternalKeyIndex int `json:"hdexternalkeyindex"`
	HdInternalKeyIndex int `json:"hdinternalkeyindex"`
}
type WalletInfoResponse struct {
	WalletVersion         int       `json:"walletversion"`
	Balance               float64   `json:"balance"`
	PrivatesendBalance    float64   `json:"privatesend_balance"`
	UnconfirmedBalance    float64   `json:"unconfirmed_balance"`
	ImmatureBalance       float64   `json:"immature_balance"`
	TxCount               int       `json:"txcount"`
	KeypoolOldest         int       `json:"keypoololdest"`
	KeypoolSize           int       `json:"keypoolsize"`
	KeypoolSizeHdInternal int       `json:"keypoolsize_hd_internal"`
	KeysLeft              int       `json:"keys_left"`
	UnlockedUntil         int       `json:"unlocked_until"`
	PayTxFee              float64   `json:"paytxfee"`
	HdChainID             string    `json:"hdchainid"`
	HdAccountCount        int       `json:"hdaccountcount"`
	HdAccount             HDAccount `json:"hdaccount"`
}
*/