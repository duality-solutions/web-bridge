package rest

import (
	"encoding/json"
	"net/http"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

func (w *WebBridgeRunner) getinfo(c *gin.Context) {
	var info models.GetInfoData
	req, _ := dynamic.NewRequest("dynamic-cli getinfo")
	err := json.Unmarshal([]byte(<-w.dynamicd.ExecCmdRequest(req)), &info)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": info})
}
