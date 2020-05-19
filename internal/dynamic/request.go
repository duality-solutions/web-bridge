package dynamic

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

var wg sync.WaitGroup
var lastID uint64 = 2

// RPCRequest is a dynamicd RPC command request
type RPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     uint64        `json:"id"`
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
			if strings.HasPrefix(c, "\"") && strings.HasSuffix(c, "\"") {
				c = c[1 : len(c)-1]
				req.Params = append(req.Params, c)
			} else if util.IsNumeric(c) {
				req.Params = append(req.Params, util.ToInt(c))
			} else {
				req.Params = append(req.Params, c)
			}

		}
	}
	wg.Add(1)
	atomic.AddUint64(&lastID, 1)
	wg.Done()
	req.ID = lastID
	return req, nil
}
