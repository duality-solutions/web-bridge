package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/gin-gonic/gin"
)

var restarting = false

//
// @Description Returns the internal web server's settings for the current running configuration
// @Accept   json
// @Produce  json
// @Success  200 {object} models.WebServerConfig	"ok"
// @Failure  500 {object} string "Internal error"
// @Router  /api/v1/config/web [get]
func (w *WebBridgeRunner) getwebserver(c *gin.Context) {
	if w.configuration != nil {
		c.JSON(http.StatusOK, gin.H{"result": w.configuration.WebServer()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}

//
// @Description Updates the internal web server's settings for the current running configuration
// @Accept   json
// @Produce  json
// @Param body body models.WebServerConfig true "WebServer"
// @Success  200 {object} models.WebServerConfig	"ok"
// @Failure  500 {object} string "Internal error"
// @Router  /api/v1/config/web [post]
func (w *WebBridgeRunner) updatewebserver(c *gin.Context) {
	if w.configuration != nil {
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
		req := models.WebServerConfig{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		if w.configuration.UpdateWebServer(req) {
			c.JSON(http.StatusOK, gin.H{"result": w.configuration.WebServer()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value. Web server configure update failed"})
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}

func waitForRestart(timeout time.Duration) {
	restarting = true
	time.Sleep(timeout)
	RestartWebServiceRouter()
	restarting = false
}

//
// @Description Restarts the internal web server using the current configuration
// @Accept   json
// @Produce  json
// @Param body body models.WebServerRestartRequest true "RestartEpoch"
// @Success  200 {object} string "ok"
// @Failure  500 {object} string "Internal error"
// @Router  /api/v1/config/web/restart [put]
func (w *WebBridgeRunner) restartwebserver(c *gin.Context) {
	if restarting {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already waiting for a restart."})
		return
	}
	if w.configuration != nil {
		var requestEpoch int64
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			strErrMsg := fmt.Sprintf("Request body read all error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		if len(body) == 0 {
			requestEpoch = 0
		} else {
			req := models.WebServerRestartRequest{}
			err = json.Unmarshal(body, &req)
			if err != nil {
				strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
				return
			}
			requestEpoch = req.RestartEpoch
		}
		currentEpoch := util.GetCurrentEpochSeconds()
		timeout := time.Second
		if requestEpoch > currentEpoch {
			timeout = time.Second * time.Duration(requestEpoch-currentEpoch)
		} else {
			timeout = time.Millisecond * 1
		}
		go waitForRestart(timeout)
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}
