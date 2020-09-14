package rest

import (
	"github.com/duality-solutions/web-bridge/init/settings"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

// WebBridgeRunner is used to run the node application.
type WebBridgeRunner struct {
	dynamicd      *dynamic.Dynamicd
	router        *gin.Engine
	configuration *settings.Configuration
	*AppShutdown
}

var runner WebBridgeRunner

// TODO: Add rate limitor
// TODO: Add custom logging
// TODO: Add bridge controller
// TODO: Add authentication
// TODO: Add RESTful API documentation with Swagger: https://github.com/swaggo/swag#getting-started

// StartWebServiceRouter is used to setup the Rest server routes
func StartWebServiceRouter(c *settings.Configuration, d *dynamic.Dynamicd, a *AppShutdown, m string) {
	gin.SetMode(m)
	runner.configuration = c
	runner.dynamicd = d
	runner.AppShutdown = a
	runner.router = gin.Default()
	api := runner.router.Group("/api")
	version := api.Group("/v1")
	setupBlockchainRoutes(version)
	setupWalletRoutes(version)
	setupBridgesRoutes(version)
	setupConfigRoutes(version)
	runner.router.Run()
}

// TODO: follow https://rest.bitcoin.com for rest endpoints
func setupBlockchainRoutes(currentVersion *gin.RouterGroup) {
	blockchain := currentVersion.Group("/blockchain")
	blockchain.POST("/jsonrpc", runner.handleJSONRPC)
	blockchain.GET("/", runner.getinfo)
	blockchain.GET("/users", runner.users)
	blockchain.GET("/users/:UserID", runner.user)
	blockchain.GET("/groups", runner.groups)
	blockchain.GET("/groups/:GroupID", runner.group)
}

// TODO: follow https://rest.bitcoin.com for rest endpoints
func setupWalletRoutes(currentVersion *gin.RouterGroup) {
	wallet := currentVersion.Group("/wallet")
	wallet.GET("/", runner.walletinfo)
	wallet.PATCH("/unlock", runner.unlockwallet)
	wallet.PATCH("/lock", runner.lockwallet)
	wallet.PATCH("/encrypt", runner.encryptwallet)
	wallet.PATCH("/changepassphrase", runner.changepassphrase)
	wallet.GET("/users", runner.walletusers)
	wallet.GET("/groups", runner.walletgroups)
	wallet.GET("/link", runner.links)
	wallet.POST("/link/request", runner.linkrequest)
	wallet.POST("/link/accept", runner.linkaccept)
	wallet.POST("/link/message", runner.sendlinkmessage)
	wallet.GET("/link/message", runner.getlinkmessages)
}

func setupBridgesRoutes(currentVersion *gin.RouterGroup) {
	bridges := currentVersion.Group("/bridges")
	bridges.GET("/", runner.bridgesinfo)
	bridges.GET("/connected", runner.connectedbridges)
	bridges.GET("/unconnected", runner.unconnectedbridges)
}

func setupConfigRoutes(currentVersion *gin.RouterGroup) {
	config := currentVersion.Group("/config")
	config.GET("/", runner.config)
	config.GET("/ice", runner.getice)
	config.PUT("/ice", runner.putice)
	config.DELETE("/ice", runner.deleteice)
	config.POST("/ice", runner.replaceice)
}
