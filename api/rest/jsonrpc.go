package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/blockchain/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

func (w *WebBridgeRunner) handleJSONRPC(c *gin.Context) {
	reqInput := models.JSONRPC{}
	err := json.NewDecoder(c.Request.Body).Decode(&reqInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	strRequest := "dynamic-cli " + reqInput.Method
	for _, param := range reqInput.Params {
		switch param.(type) {
		case int:
			val, ok := param.(int)
			if ok {
				strRequest += ` ` + fmt.Sprintf("%v", val)
			}
		case float64:
			val, ok := param.(float64)
			if ok {
				strRequest += ` ` + fmt.Sprintf("%v", val)
			}
		case bool:
			val, ok := param.(bool)
			if ok {
				strRequest += ` ` + fmt.Sprintf("%v", val)
			}
		case string:
			strRequest += ` "` + fmt.Sprintf("%v", param) + `"`
		}
	}
	reqOutput, _ := dynamic.NewRequest(strRequest)
	response, _ := <-w.dynamicd.ExecCmdRequest(reqOutput)
	var result interface{}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Result unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}
