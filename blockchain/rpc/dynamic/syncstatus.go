package dynamic

import (
	"encoding/json"
)

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

// GetSyncStatus returns the dynamicd blockchain status
func (d *Dynamicd) GetSyncStatus() (*SyncStatus, error) {
	var status SyncStatus
	req, _ := NewRequest("dynamic-cli syncstatus")
	err := json.Unmarshal([]byte(<-d.ExecCmdRequest(req)), &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetNumberOfConnections returns the number of active blockchain peer connections
func (d *Dynamicd) GetNumberOfConnections() (int, error) {
	var status SyncStatus
	req, _ := NewRequest("dynamic-cli syncstatus")
	err := json.Unmarshal([]byte(<-d.ExecCmdRequest(req)), &status)
	if err != nil {
		return 0, err
	}
	return status.Peers, nil
}
