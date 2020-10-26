package models

// TransactionsResponse contains the wallet's send and received transaction list
// swagger:parameters models.TransactionsResponse
type TransactionsResponse struct {
	Account           string   `json:"account"`
	Address           string   `json:"address"`
	Category          string   `json:"Category"`
	Amount            float64  `json:"Amount"`
	Label             string   `json:"label"`
	VOut              int      `json:"vout"`
	Confirmations     int      `json:"confirmations"`
	InstantLock       bool     `json:"instantlock"`
	BlockHash         string   `json:"blockhash"`
	BlockIndex        int      `json:"blockindex"`
	BlockTime         int      `json:"blocktime"`
	TxID              string   `json:"txid"`
	WalletConflicts   []string `json:"walletconflicts"`
	Time              int64    `json:"time"`
	TimeReceived      int64    `json:"timereceived"`
	BIP125Replaceable string   `json:"bip125-replaceable"`
}
