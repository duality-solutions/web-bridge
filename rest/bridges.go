package rest

import (
	"fmt"
	"net/http"

	"github.com/duality-solutions/web-bridge/bridge"
	"github.com/gin-gonic/gin"
)

type bridgeInfo struct {
	SessionID          uint16 `json:"session_id"`
	State              string `json:"state"`
	MyAccount          string `json:"my_account"`
	LinkAccount        string `json:"link_account"`
	LinkID             string `json:"link_id"`
	OnOpenEpoch        int64  `json:"on_open_epoch"`
	OnStateChangeEpoch int64  `json:"on_state_changed_epoch"`
	OnLastDataEpoch    int64  `json:"on_last_data_epoch"`
	OnErrorEpoch       int64  `json:"on_error_epoch"`
	RTCState           string `json:"rtc_status"`
	ListenPort         uint16 `json:"listen_port"`
}

func (w *WebBridgeRunner) bridgesinfo(c *gin.Context) {
	controler, err := bridge.Controler()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	bridges := controler.AllBridges()
	var ret []bridgeInfo = make([]bridgeInfo, len(bridges))
	i := 0
	for _, bridge := range bridges {
		var info = bridgeInfo{
			SessionID:          bridge.SessionID,
			State:              bridge.State().String(),
			MyAccount:          bridge.MyAccount,
			LinkAccount:        bridge.LinkAccount,
			OnOpenEpoch:        bridge.OnOpenEpoch(),
			OnLastDataEpoch:    bridge.OnLastDataEpoch(),
			OnErrorEpoch:       bridge.OnErrorEpoch(),
			OnStateChangeEpoch: bridge.OnStateChangeEpoch(),
			RTCState:           bridge.RTCState(),
			ListenPort:         bridge.ListenPort(),
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
	var ret []bridgeInfo = make([]bridgeInfo, len(bridges))
	for _, bridge := range bridges {
		var info = bridgeInfo{
			SessionID:          bridge.SessionID,
			State:              bridge.State().String(),
			MyAccount:          bridge.MyAccount,
			LinkAccount:        bridge.LinkAccount,
			OnOpenEpoch:        bridge.OnOpenEpoch(),
			OnLastDataEpoch:    bridge.OnLastDataEpoch(),
			OnErrorEpoch:       bridge.OnErrorEpoch(),
			OnStateChangeEpoch: bridge.OnStateChangeEpoch(),
			RTCState:           bridge.RTCState(),
			ListenPort:         bridge.ListenPort(),
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
	var ret []bridgeInfo = make([]bridgeInfo, len(bridges))
	for _, bridge := range bridges {
		var info = bridgeInfo{
			SessionID:          bridge.SessionID,
			State:              bridge.State().String(),
			MyAccount:          bridge.MyAccount,
			LinkAccount:        bridge.LinkAccount,
			OnOpenEpoch:        bridge.OnOpenEpoch(),
			OnLastDataEpoch:    bridge.OnLastDataEpoch(),
			OnErrorEpoch:       bridge.OnErrorEpoch(),
			OnStateChangeEpoch: bridge.OnStateChangeEpoch(),
			RTCState:           bridge.RTCState(),
			ListenPort:         bridge.ListenPort(),
		}
		ret[i] = info
		i++
	}
	c.JSON(http.StatusOK, gin.H{"result": ret})
}
