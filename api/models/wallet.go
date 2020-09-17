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

type WalletSeed struct {
	HDSeed             string `json:"hdseed"`
	Mnemonic           string `json:"mnemonic"`
	MnemonicPassphrase string `json:"mnemonicpassphrase"`
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
