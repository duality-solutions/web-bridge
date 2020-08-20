package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

type unlockWalletRequest struct {
	Passphrase string `json:"passphrase"`
	Timeout    int    `json:"timeout"`
	MixingOnly bool   `json:"mixingonly"`
}

func (w *WebBridgeRunner) unlockwallet(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body read all error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty."})
		return
	}
	req := unlockWalletRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(req.Passphrase) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passphrase can not be empty."})
		return
	}
	if req.Timeout == 0 {
		req.Timeout = 60000
	}
	strRequest := `dynamic-cli walletpassphrase "` + req.Passphrase + `"` + ` ` + fmt.Sprintf("%v", req.Timeout)
	if req.MixingOnly == true {
		strRequest += ` 1`
	}
	strCommand, _ := dynamic.NewRequest(strRequest)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	if strings.HasPrefix(response, "null") {
		c.JSON(http.StatusOK, gin.H{"result": "successful"})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": response})
	}
}

func (w *WebBridgeRunner) lockwallet(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli walletlock`)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	if strings.HasPrefix(response, "null") {
		c.JSON(http.StatusOK, gin.H{"result": "successful"})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": response})
	}
}

// encryptwallet
// walletlock
