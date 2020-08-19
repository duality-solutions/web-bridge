package rest

import (
	"encoding/json"
	"net/http"

	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

type getInfoData struct {
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

func (w *WebBridgeRunner) getinfo(c *gin.Context) {
	var info getInfoData
	req, _ := dynamic.NewRequest("dynamic-cli getinfo")
	err := json.Unmarshal([]byte(<-w.dynamicd.ExecCmdRequest(req)), &info)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": info})
}
