package webbridge

// SyncStatus is a dynamicd RPC command result
type SyncStatus struct {
	ChainName          string  `json:"chain_name"`
	Version            int     `json:"version"`
	Peers              int     `json:"peers"`
	Headers            int     `json:"headers"`
	Blocks             int     `json:"blocks"`
	CurrentBlockHeight int     `json:"current_block_height"`
	SyncProgress       float32 `json:"sync_progress"`
	StatusDescription  string  `json:"status_description"`
	FullySynced        bool    `json:"fully_synced"`
	Failed             bool    `json:"failed"`
}