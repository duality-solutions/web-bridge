package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

/*
getusers "search string" "records per page" "page returned"

Arguments:
1. search string        (string, optional)  Search for userid
2. records per page     (int, optional)  If paging, the number of records per page
3. page returned        (int, optional)  If paging, the page number to return

Lists all BDAP user accounts in the "public" OU for the "bdap.io" domain.

Result:
{(json objects)
  "common_name"             (string)  Account common name
  "object_full_path"        (string)  Account fully qualified domain name (FQDN)
  "wallet_address"          (string)  Account wallet address
  "dht_publickey"           (string)  Account DHT public key
  }

Examples
> dynamic-cli getusers

As a JSON-RPC call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getusers", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:33350/
 (code -1)
*/
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
