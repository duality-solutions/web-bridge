package models

type SyncStatus struct {
	Blocks             int     `json:"blocks"`
	ChainName          string  `json:"chain_name"`
	CurrentBlockHeight int     `json:"current_block_height"`
	Failed             bool    `json:"failed"`
	FullySynced        bool    `json:"fully_synced"`
	Headers            int     `json:"headers"`
	Peers              int     `json:"peers"`
	StatusDescription  string  `json:"status_description"`
	SyncProgress       float64 `json:"sync_progress"`
	ClientVersion      int     `json:"version"`
}
