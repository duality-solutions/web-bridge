package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w *WebBridgeRunner) shutdown(c *gin.Context) {
	if w.shutdownApp != nil {
		w.shutdownApp.ShutdownAppliction()
		c.JSON(http.StatusOK, gin.H{"result": "WebBridge is shutting down."})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable shutdownApp is null."})
	}
}
