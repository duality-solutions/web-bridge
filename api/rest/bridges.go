package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/bridge"
	"github.com/gin-gonic/gin"
)

func (w *WebBridgeRunner) bridgesinfo(c *gin.Context) {
	controler, err := bridge.Controler()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	bridges := controler.AllBridges()
	var ret []models.BridgeInfo = make([]models.BridgeInfo, len(bridges))
	i := 0
	for _, bridge := range bridges {
		var info = models.BridgeInfo{
			SessionID:          bridge.SessionID,
			LinkID:             bridge.LinkID(),
			State:              bridge.State().String(),
			MyAccount:          bridge.MyAccount,
			LinkAccount:        bridge.LinkAccount,
			OnOpenEpoch:        bridge.OnOpenEpoch(),
			OnLastDataEpoch:    bridge.OnLastDataEpoch(),
			OnErrorEpoch:       bridge.OnErrorEpoch(),
			OnStateChangeEpoch: bridge.OnStateChangeEpoch(),
			RTCState:           bridge.RTCState(),
			HTTPListenPort:     bridge.ListenPort(),
			HTTPSListenPort:    bridge.ListenPort() + 1,
		}
		ret[i] = info
		i++
	}
	c.JSON(http.StatusOK, gin.H{"result": ret})
}

func (w *WebBridgeRunner) connectedbridges(c *gin.Context) {
	controler, err := bridge.Controler()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	i := 0
	bridges := controler.Connected()
	var ret []models.BridgeInfo = make([]models.BridgeInfo, len(bridges))
	for _, bridge := range bridges {
		var info = models.BridgeInfo{
			SessionID:          bridge.SessionID,
			LinkID:             bridge.LinkID(),
			State:              bridge.State().String(),
			MyAccount:          bridge.MyAccount,
			LinkAccount:        bridge.LinkAccount,
			OnOpenEpoch:        bridge.OnOpenEpoch(),
			OnLastDataEpoch:    bridge.OnLastDataEpoch(),
			OnErrorEpoch:       bridge.OnErrorEpoch(),
			OnStateChangeEpoch: bridge.OnStateChangeEpoch(),
			RTCState:           bridge.RTCState(),
			HTTPListenPort:     bridge.ListenPort(),
			HTTPSListenPort:    bridge.ListenPort() + 1,
		}
		ret[i] = info
		i++
	}
	c.JSON(http.StatusOK, gin.H{"result": ret})
}

func (w *WebBridgeRunner) unconnectedbridges(c *gin.Context) {
	controler, err := bridge.Controler()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	i := 0
	bridges := controler.Unconnected()
	var ret []models.BridgeInfo = make([]models.BridgeInfo, len(bridges))
	for _, bridge := range bridges {
		var info = models.BridgeInfo{
			SessionID:          bridge.SessionID,
			LinkID:             bridge.LinkID(),
			State:              bridge.State().String(),
			MyAccount:          bridge.MyAccount,
			LinkAccount:        bridge.LinkAccount,
			OnOpenEpoch:        bridge.OnOpenEpoch(),
			OnLastDataEpoch:    bridge.OnLastDataEpoch(),
			OnErrorEpoch:       bridge.OnErrorEpoch(),
			OnStateChangeEpoch: bridge.OnStateChangeEpoch(),
			RTCState:           bridge.RTCState(),
			HTTPListenPort:     bridge.ListenPort(),
			HTTPSListenPort:    bridge.ListenPort() + 1,
		}
		ret[i] = info
		i++
	}
	c.JSON(http.StatusOK, gin.H{"result": ret})
}

func shutdownBridge(b *bridge.Bridge) (int, error) {
	var status = http.StatusOK
	controler, err := bridge.Controler()
	if err != nil {
		status = http.StatusInternalServerError
		return status, fmt.Errorf("Bridge controller error %v", err)
	}
	conn := controler.GetConnected(b.LinkID())
	if conn != nil {
		b.ShutdownHTTPProxyServers()
		controler.MoveConnectedToUnconnected(b)
	}
	return status, nil
}

func restartBridge(b *bridge.Bridge) (int, error) {
	// TODO: send VGP message
	var status = http.StatusOK
	controler, err := bridge.Controler()
	if err != nil {
		status = http.StatusInternalServerError
		return status, fmt.Errorf("Bridge controller error %v", err)
	}
	newBridge := bridge.ResetBridge(b)
	controler.PutUnconnected(newBridge)
	return status, nil
}

func getBridgeFromBody(reqBody io.ReadCloser) (*bridge.Bridge, int, error) {
	var status = http.StatusBadRequest
	body, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil, status, fmt.Errorf("Request body read all error %v", err)
	}
	if len(body) == 0 {
		return nil, status, fmt.Errorf("Body is empty")
	}
	req := models.BridgeRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, status, fmt.Errorf("Request body JSON unmarshal error %v", err)
	}
	if len(req.LinkID) == 0 {
		return nil, status, fmt.Errorf("Request body contains an empty link ID")
	}
	controler, err := bridge.Controler()
	if err != nil {
		status = http.StatusInternalServerError
		return nil, status, fmt.Errorf("Bridge controller error %v", err)
	}
	b := controler.AllBridges()[req.LinkID]
	if b != nil {
		status = http.StatusOK
	}
	return b, status, nil
}

//
// @Description Restarts the specified bridge
// @Accept  json
// @Produce  json
// @Param body body models.BridgeRequest true "Bridge"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal error"
// @Router /api/v1/bridges/restart [post]
func (w *WebBridgeRunner) restartbridge(c *gin.Context) {
	b, status, err := getBridgeFromBody(c.Request.Body)
	if err != nil {
		c.JSON(status, gin.H{"error": err})
		return
	}
	status, err = shutdownBridge(b)
	if err != nil {
		c.JSON(status, gin.H{"error": err})
		return
	}
	status, err = restartBridge(b)
	if err != nil {
		c.JSON(status, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": b})
}

func (w *WebBridgeRunner) startbridge(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"result": "stub"})
}

func (w *WebBridgeRunner) stopbridge(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"result": "stub"})
}
