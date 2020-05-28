package dynamic

import (
	"encoding/json"
)

// GetInfo returns an object containing various state info.
type GetInfo struct {
	Balance            float64 `json:"balance"`
	Blocks             int     `json:"blocks"`
	Connections        int     `json:"connections"`
	Difficulty         float64 `json:"difficulty"`
	Errors             string  `json:"errors"`
	KeyPoolOldest      int     `json:"keypoololdest"`
	KeyPoolSize        int     `json:"keypoolsize"`
	PayTxFee           float64 `json:"paytxfee"`
	PrivateSendBalance float64 `json:"privatesend_balance"`
	ProtocolVersion    int     `json:"protocolversion"`
	Proxy              string  `json:"proxy"`
	RelayFee           float64 `json:"relayfee"`
	Testnet            bool    `json:"testnet"`
	TimeOffset         int     `json:"timeoffset"`
	UnlockedUntil      int     `json:"unlocked_until"`
	Version            int     `json:"version"`
	Walletversion      int     `json:"walletversion"`
}

// GetInfo returns all BDAP accounts from the wallet
func (d *Dynamicd) GetInfo() (*GetInfo, error) {
	var info GetInfo
	req, _ := NewRequest("dynamic-cli getinfo")
	err := json.Unmarshal([]byte(<-d.ExecCmdRequest(req)), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
