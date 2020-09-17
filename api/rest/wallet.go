package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

func (w *WebBridgeRunner) unlockwallet(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body read all error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}
	req := models.UnlockWalletRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(req.Passphrase) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passphrase can not be empty"})
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
		var result interface{}
		err = json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		strResult := fmt.Sprintf("%v", result)
		if strings.Contains(strResult, "Wallet is already fully unlocked") {
			c.JSON(http.StatusOK, gin.H{"result": "Wallet is already fully unlocked"})
			return
		}
		if strings.Contains(strResult, "running with an unencrypted wallet") {
			c.JSON(http.StatusOK, gin.H{"result": "Can not unlock an unencrypted wallet"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (w *WebBridgeRunner) lockwallet(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli walletlock`)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	if strings.HasPrefix(response, "null") {
		c.JSON(http.StatusOK, gin.H{"result": "successful"})
	} else {
		var result interface{}
		err := json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		strResult := fmt.Sprintf("%v", result)
		if strings.Contains(strResult, "Wallet is already fully locked") {
			c.JSON(http.StatusOK, gin.H{"result": "Wallet is already fully locked"})
			return
		}
		if strings.Contains(strResult, "running with an unencrypted wallet") {
			c.JSON(http.StatusOK, gin.H{"result": "Can not lock an unencrypted wallet"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (w *WebBridgeRunner) encryptwallet(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body read all error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}
	req := models.EncryptWalletRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(req.Passphrase) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passphrase can not be empty"})
		return
	}
	strRequest := `dynamic-cli encryptwallet "` + req.Passphrase + `"`
	strCommand, _ := dynamic.NewRequest(strRequest)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	if strings.HasPrefix(response, "null") {
		c.JSON(http.StatusOK, gin.H{"result": "successful"})
	} else {
		var result interface{}
		err = json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		strResult := fmt.Sprintf("%v", result)
		if strings.Contains(strResult, "running with an encrypted wallet") {
			c.JSON(http.StatusOK, gin.H{"result": "Can not encrypt a wallet that is already encrypted"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (w *WebBridgeRunner) changepassphrase(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body read all error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}
	req := models.ChangePassphraseRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(req.OldPassphrase) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old passphrase can not be empty"})
		return
	}
	if len(req.NewPassphrase) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New passphrase can not be empty"})
		return
	}
	strRequest := `dynamic-cli walletpassphrasechange "` + req.OldPassphrase + `"` + ` "` + req.NewPassphrase + `"`
	strCommand, _ := dynamic.NewRequest(strRequest)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	if strings.HasPrefix(response, "null") {
		c.JSON(http.StatusOK, gin.H{"result": "successful"})
	} else {
		var result interface{}
		err = json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		strResult := fmt.Sprintf("%v", result)
		if strings.Contains(strResult, "running with an unencrypted wallet") {
			c.JSON(http.StatusOK, gin.H{"error": "Can not change passphase when the wallet is not encrypted"})
			return
		}
		if strings.Contains(strResult, "wallet passphrase entered was incorrect") {
			c.JSON(http.StatusOK, gin.H{"error": "Incorrect wallet passphrase"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (w *WebBridgeRunner) walletinfo(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli getwalletinfo`)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	var result interface{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (w *WebBridgeRunner) mnemonic(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli dumphdinfo`)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	if strings.Contains(response, "Please enter the wallet passphrase with walletpassphrase first") {
		result := models.RPCError{}
		err := json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("Response JSON unmarshal error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	result := models.WalletSeed{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}

// getwalletinfo also add if locked or not.
