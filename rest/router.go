package rest

import (
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

// WebBridgeRunner is used to run the node application.
type WebBridgeRunner struct {
	dynamicd *dynamic.Dynamicd
	router   *gin.Engine
}

var runner WebBridgeRunner

// TODO: Add rate limitor
// TODO: Add custom logging
// TODO: Add bridge controller
// TODO: Add authentication
// TODO: Add RESTful API documentation with Swagger: https://github.com/swaggo/swag#getting-started

// StartWebServiceRouter is used to setup the Rest server routes
func StartWebServiceRouter(dynamicd *dynamic.Dynamicd, mode string) {
	gin.SetMode(mode)
	runner.dynamicd = dynamicd
	runner.router = gin.Default()
	api := runner.router.Group("/api")
	version := api.Group("/v1")
	setupBlockchainRoutes(version)
	setupWalletRoutes(version)
	runner.router.Run()
}

// TODO: follow https://rest.bitcoin.com for rest endpoints
func setupBlockchainRoutes(currentVersion *gin.RouterGroup) {
	blockchain := currentVersion.Group("/blockchain")
	blockchain.POST("/jsonrpc", runner.handleJSONRPC)
	blockchain.GET("/info", runner.getinfo)
	blockchain.GET("/users", runner.users)
	blockchain.GET("/users/:UserID", runner.user)
	blockchain.GET("/groups", runner.groups)
	blockchain.GET("/groups/:GroupID", runner.group)
}

// TODO: follow https://rest.bitcoin.com for rest endpoints
func setupWalletRoutes(currentVersion *gin.RouterGroup) {
	wallet := currentVersion.Group("/wallet")
	wallet.GET("/info", runner.walletinfo)
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
