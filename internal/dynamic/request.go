package dynamic

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

var wg sync.WaitGroup
var lastID int64 = 2

// RPCRequest is a type for raw JSON-RPC 1.0 requests.  The Method field identifies
// the specific command type which in turns leads to different parameters.
// Callers typically will not use this directly since this package provides a
// statically typed command infrastructure which handles creation of these
// requests, however this struct it being exported in case the caller wants to
// construct raw requests for some reason.
type RPCRequest struct {
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int64         `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
}

// NewRequest returns a new JSON-RPC 1.0 request object given the raw command.
func NewRequest(cmd string) (*RPCRequest, error) {
	var req RPCRequest
	if len(cmd) < 12 {
		return nil, fmt.Errorf("incorrect cmd size %s", cmd)
	}
	if !strings.HasPrefix(cmd, "dynamic-cli") {
		return nil, fmt.Errorf("incorrect prefix %s", cmd)
	}
	cmd = strings.Replace(cmd, "dynamic-cli ", "", -1)
	i := strings.Index(cmd, " ")
	if i > 0 {
		req.Method = strings.TrimSpace(cmd[:i])
		cmd = cmd[i+1:]
		err := req.parseCmd(cmd)
		if err != nil {
			return nil, err
		}
	} else {
		req.Method = strings.TrimSpace(cmd)
	}
	req.JSONRPC = "1.0"
	wg.Add(1)
	atomic.AddInt64(&lastID, 1)
	wg.Done()
	req.ID = lastID
	return &req, nil
}

func (req *RPCRequest) parseCmd(cmd string) error {
	// Split cmd string
	r := NewCmdReader(strings.NewReader(cmd))
	r.Delimiter = ' ' // space
	paramsRaw, err := r.Read()
	if err != nil {
		return err
	}
	for _, val := range paramsRaw {
		if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
			val = val[1 : len(val)-1]
			req.Params = append(req.Params, val)
		} else if util.IsNumeric(val) {
			req.Params = append(req.Params, util.ToInt(val))
		} else {
			req.Params = append(req.Params, val)
		}
	}
	return nil
}
