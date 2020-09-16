package rest

import (
	"fmt"
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
