package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/gin-gonic/gin"
)

//
// @Description Returns all the ICE servers in current running configuration
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.IceServerConfig	"ok"
// @Failure 500 {object} string "Internal error"
// @Router /api/v1/ice [get]
func (w *WebBridgeRunner) getice(c *gin.Context) {
	if w.configuration != nil {
		c.JSON(http.StatusOK, w.configuration.IceServers())
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}

func (w *WebBridgeRunner) findIceSetting(iceSetting models.IceServerConfig) (int, error) {
	if w.configuration == nil {
		return -1, fmt.Errorf("configuration variable is null")
	}
	iceServers := *w.configuration.IceServers()
	for i, ice := range iceServers {
		if ice.URL == iceSetting.URL {
			return i, nil
		}
	}
	return -1, fmt.Errorf("ICE settings not found")
}

//
// @Description Add or update an ICE server in current configuration and saves the changes to file
// @Accept  json
// @Produce json
// @Param body body models.IceServerConfig true "ICE Configuration"
// @Success 200 {object} []models.IceServerConfig	"ok"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {object} string "Internal error"
// @Router /api/v1/ice [put]
func (w *WebBridgeRunner) putice(c *gin.Context) {
	if w.configuration != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			strErrMsg := fmt.Sprintf("Request body read all error %v. Can not add ICE server to configuration.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		if len(body) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty. Can not add ICE server to configuration."})
			return
		}
		req := models.IceServerConfig{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v. Can not add ICE server to configuration.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		if len(req.URL) == 0 {
			strErrMsg := fmt.Sprintf("URL is empty. Can not add ICE server to configuration.")
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		index, err := w.findIceSetting(req)
		if err != nil {
			w.configuration.AddIceServer(req)
		} else {
			w.configuration.UpdateIceServer(index, req)
		}
		c.JSON(http.StatusOK, w.configuration.IceServers())
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null. Can not add ICE server to configuration."})
	}
}

//
// @Description Delete an ICE server in current configuration and saves the changes to file
// @Accept  json
// @Produce json
// @Param body body models.IceServerConfig true "ICE Configuration"
// @Success 200 {object} []models.IceServerConfig	"ok"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal error"
// @Router /api/v1/ice [delete]
func (w *WebBridgeRunner) deleteice(c *gin.Context) {
	if w.configuration != nil {
		if len(*w.configuration.IceServers()) < 2 {
			strErrMsg := fmt.Sprintf("Can not delete last ICE server in configuration")
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			strErrMsg := fmt.Sprintf("Request body read all error %v. Can not detele from ICE server configuration list.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		if len(body) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty. Can not detele from ICE server configuration list."})
			return
		}
		req := models.IceServerConfig{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v. Can not detele from ICE server configuration list.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		if len(req.URL) == 0 {
			strErrMsg := fmt.Sprintf("URL is empty. Can not detele from ICE server configuration list.")
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		index, err := w.findIceSetting(req)
		if err == nil {
			w.configuration.DeleteIceServer(index)
		} else {
			strErrMsg := fmt.Sprintf("Setting not found by URL %v. Can not detele from ICE server configuration list.", req.URL)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		}
		c.JSON(http.StatusOK, w.configuration.IceServers())
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null. Can not detele from ICE server configuration list."})
	}
}

//
// @Description Replaces an ICE server in current configuration and saves the changes to file
// @Accept  json
// @Produce  json
// @Param body body []models.IceServerConfig true "ICE Configuration File"
// @Success 200 {object} []models.IceServerConfig	"ok"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal error"
// @Router /api/v1/ice [post]
func (w *WebBridgeRunner) replaceice(c *gin.Context) {
	if w.configuration != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			strErrMsg := fmt.Sprintf("Request body read all error %v. Can not detele from ICE server configuration list.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		if len(body) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty. Can not detele from ICE server configuration list."})
			return
		}
		req := []models.IceServerConfig{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v. Can not detele from ICE server configuration list.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		if len(req) == 0 {
			strErrMsg := fmt.Sprintf("URL is empty. Can not detele from ICE server configuration list.")
			c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
			return
		}
		w.configuration.ReplaceIceServers(req)
		c.JSON(http.StatusOK, w.configuration.IceServers())
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null. Can not detele from ICE server configuration list."})
	}
}
