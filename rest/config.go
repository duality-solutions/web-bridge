package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/duality-solutions/web-bridge/init/settings"
	"github.com/gin-gonic/gin"
)

func (w *WebBridgeRunner) config(c *gin.Context) {
	if configuration != nil {
		c.JSON(http.StatusOK, gin.H{"result": configuration})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}

func (w *WebBridgeRunner) getice(c *gin.Context) {
	if configuration != nil {
		c.JSON(http.StatusOK, gin.H{"result": configuration.IceServers})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}

func findIceSetting(iceSetting settings.IceServerConfig) (int, error) {
	if configuration == nil {
		return -1, fmt.Errorf("configuration variable is null")
	}
	for i, ice := range configuration.IceServers {
		if ice.URL == iceSetting.URL {
			return i, nil
		}
	}
	return -1, fmt.Errorf("ICE settings not found")
}

// TODO configuration.IceServers needs a lock
// TODO write changes to config file while locked
func (w *WebBridgeRunner) putice(c *gin.Context) {
	if configuration != nil {
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
		req := settings.IceServerConfig{}
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
		index, err := findIceSetting(req)
		if err != nil {
			configuration.IceServers = append(configuration.IceServers, req)
		} else {
			configuration.IceServers[index] = req
		}
		c.JSON(http.StatusOK, gin.H{"result": configuration.IceServers})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}
