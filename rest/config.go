package rest

import (
	"net/http"
	"time"

	"github.com/duality-solutions/web-bridge/bridge"
	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

// AppShutdown stores the channels and commands needed to safely shutdown the WebBridge service
type AppShutdown struct {
	Close       *chan struct{}
	StopWatcher *chan struct{}
	StopBridges *chan struct{}
	Dynamicd    *dynamic.Dynamicd
}

// ShutdownAppliction safely shuts down bridge proxies, process watcher, and the Dynamic daemon before exiting
func (a *AppShutdown) ShutdownAppliction() {
	close(*a.StopWatcher)
	bridge.ShutdownBridges(a.StopBridges)
	// Stop dynamicd
	reqStop, _ := dynamic.NewRequest("dynamic-cli stop")
	respStop, _ := util.BeautifyJSON(<-a.Dynamicd.ExecCmdRequest(reqStop))
	util.Info.Println(respStop)
	time.Sleep(time.Second * 5)
	util.Info.Println("Looking for dynamicd process pid", a.Dynamicd.Cmd.Process.Pid)
	util.WaitForProcessShutdown(a.Dynamicd.Cmd.Process, time.Second*240)
	util.Info.Println("Exit.")
	util.EndDebugLogFile(30)
	close(*a.Close)
}

func (w *WebBridgeRunner) config(c *gin.Context) {
	if w.configuration != nil {
		c.JSON(http.StatusOK, gin.H{"result": w.configuration.ToJSON()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}
