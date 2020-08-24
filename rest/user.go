package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

func (w *WebBridgeRunner) users(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli getusers`)
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

func (w *WebBridgeRunner) user(c *gin.Context) {
	userID := c.Param("UserID")
	cmd := `dynamic-cli getuserinfo "` + userID + `"`
	strCommand, err := dynamic.NewRequest(cmd)
	if err != nil {
		strErrMsg := fmt.Sprintf("NewRequest cmd(%v), error: %v", cmd, err)
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

type account struct {
	OID                string `json:"oid"`
	Version            int    `json:"version"`
	DomainComponent    string `json:"domain_component"`
	CommonName         string `json:"common_name"`
	OrganizationalUnit string `json:"organizational_unit"`
	OrganizationName   string `json:"organization_name"`
	ObjectID           string `json:"object_id"`
	ObjectFullPath     string `json:"object_full_path"`
	ObjectType         string `json:"object_type"`
	WalletAddress      string `json:"wallet_address"`
	Public             int    `json:"public"`
	DhtPublickey       string `json:"dht_publickey"`
	LinkAddress        string `json:"link_address"`
	TxID               string `json:"txid"`
	Time               int    `json:"time"`
	ExpiresOn          int    `json:"expires_on"`
	Expired            bool   `json:"expired"`
}

func (w *WebBridgeRunner) walletusers(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli mybdapaccounts`)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	myAccounts := make(map[string]account)
	err := json.Unmarshal([]byte(response), &myAccounts)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}

	myUsers := make(map[string]account)
	for i, account := range myAccounts {
		if account.ObjectType == "User Entry" {
			myUsers[i] = account
		}
	}

	c.JSON(http.StatusOK, gin.H{"result": myUsers})
}
