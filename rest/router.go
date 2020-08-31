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

func StartWebServiceRouter(dynamicd *dynamic.Dynamicd, mode string) {
	gin.SetMode(mode)
	runner.dynamicd = dynamicd
	runner.router = gin.Default()
	setupBlockchainRoutes()
	runner.router.Run()
}

// TODO: follow https://rest.bitcoin.com for rest endpoints
func setupBlockchainRoutes() {
	api := runner.router.Group("/api")
	v1 := api.Group("/v1")
	blockchain := v1.Group("/blockchain")
	blockchain.POST("/jsonrpc", runner.handleJSONRPC)
	blockchain.GET("/info", runner.getinfo)
	blockchain.GET("/wallet/info", runner.walletinfo)
	blockchain.PATCH("/wallet/unlock", runner.unlockwallet)
	blockchain.PATCH("/wallet/lock", runner.lockwallet)
	blockchain.PATCH("/wallet/encrypt", runner.encryptwallet)
	blockchain.PATCH("/wallet/changepassphrase", runner.changepassphrase)
	blockchain.GET("/wallet/users", runner.walletusers)
	blockchain.GET("/wallet/groups", runner.walletgroups)
	blockchain.GET("/wallet/link", runner.links)
	blockchain.POST("/wallet/link/request", runner.linkrequest)
	blockchain.POST("/wallet/link/accept", runner.linkaccept)
	blockchain.GET("/users", runner.users)
	blockchain.GET("/users/:UserID", runner.user)
	blockchain.GET("/groups", runner.groups)
	blockchain.GET("/groups/:GroupID", runner.group)
}
