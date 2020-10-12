package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w *WebBridgeRunner) config(c *gin.Context) {
	if w.configuration != nil {
		c.JSON(http.StatusOK, w.configuration.ToJSON())
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable is null."})
	}
}
