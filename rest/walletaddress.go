package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

type listAddressReponse struct {
	WalletAddress string  `json:"walletaddress"`
	Balance       float32 `json:"balance"`
}

// address/details
func (w *WebBridgeRunner) listaddresses(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli listaddressbalances`)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	result := listAddressReponse{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
	return
}

// address/details/:Address
func (w *WebBridgeRunner) addressdetails(c *gin.Context) {
	address := c.Param("Address")
	cmd := `dynamic-cli getaddressdeltas '{"addresses":["` + address + `"]}'`
	c.JSON(http.StatusOK, gin.H{"cmd": cmd})
	strCommand, err := dynamic.NewRequest(cmd)
	if err != nil {
		strErrMsg := fmt.Sprintf("NewRequest error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	var result interface{}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
	return
}
