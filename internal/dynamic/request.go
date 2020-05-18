package dynamic

import (
	"fmt"
	"strings"
)

var lastID int64 = 2

// RPCRequest is a dynamicd RPC command request
type RPCRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int64    `json:"id"`
}

// NewRequest return as new RPCRequest struct from the given cmd text
func NewRequest(cmd string) (RPCRequest, error) {
	var req RPCRequest
	if len(cmd) < 12 {
		return req, fmt.Errorf("incorrect cmd size %s", cmd)
	}
	if !strings.HasPrefix(cmd, "dynamic-cli") {
		return req, fmt.Errorf("incorrect prefix %s", cmd)
	}
	cmd = strings.Replace(cmd, "dynamic-cli ", "", -1)
	for i, c := range strings.Split(cmd, " ") {
		if i == 0 {
			req.Method = c
		} else {
			req.Params = append(req.Params, c)
		}
	}
	lastID++
	req.ID = lastID
	return req, nil
}
