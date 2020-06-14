package dynamic

import (
	"encoding/json"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

// WaitForWalletCreated waits until the Dynamic wallet is loaded
func (d *Dynamicd) WaitForWalletCreated() {
	walletMissing := true
	for walletMissing {
		req, _ := NewRequest("dynamic-cli getinfo")
		res := <-d.ExecCmdRequest(req)
		var info GetInfo
		err := json.Unmarshal([]byte(<-d.ExecCmdRequest(req)), &info)
		if err == nil {
			walletMissing = false
		} else {
			var rpcError RPCErrorResponse
			err = json.Unmarshal([]byte(res), &rpcError)
			if err == nil {
				util.Error.Println("WaitForWalletCreated ...", rpcError.Error.Message)
			}
		}
	}
}
