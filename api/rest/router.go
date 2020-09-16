package rest

import (
	"github.com/duality-solutions/web-bridge/configs/settings"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"

	_ "github.com/duality-solutions/web-bridge/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// WebBridgeRunner is used to run the node application.
type WebBridgeRunner struct {
	dynamicd      *dynamic.Dynamicd
	router        *gin.Engine
	configuration *settings.Configuration
	shutdownApp   *AppShutdown
}

var runner WebBridgeRunner

// TODO: Add rate limitor
// TODO: Add custom logging
// TODO: Add authentication

// StartWebServiceRouter is used to setup the Rest server routes
func StartWebServiceRouter(c *settings.Configuration, d *dynamic.Dynamicd, a *AppShutdown, m string) {
	gin.SetMode(m)
	runner.configuration = c
	runner.dynamicd = d
	runner.shutdownApp = a
	runner.router = gin.Default()
	api := runner.router.Group("/api")
	version := api.Group("/v1")
	version.POST("/shutdown", runner.shutdown)
	setupBlockchainRoutes(version)
	setupWalletRoutes(version)
	setupBridgesRoutes(version)
	setupConfigRoutes(version)
	setupSwagger(runner.router)
	runner.router.Run()
}

// @title WebBridge Restful API Documentation
// @version 1.0
// @description WebBridge Rest API discovery website.
// @termsOfService http://www.duality.solutions/webbridge/terms

// @contact.name API Support
// @contact.url http://www.duality.solutions/support
// @contact.email support@duality.solutions

// @license.name Duality
// @license.url https://github.com/duality-solutions/web-bridge/blob/master/LICENSE.md

// @host http://docs.webbridge.duality.solutions
// @BasePath /api/v1
func setupSwagger(root *gin.Engine) {
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// TODO: follow https://rest.bitcoin.com for rest endpoints
func setupBlockchainRoutes(currentVersion *gin.RouterGroup) {
	blockchain := currentVersion.Group("/blockchain")
	blockchain.POST("/jsonrpc", runner.handleJSONRPC)
	blockchain.GET("/", runner.getinfo)
	blockchain.GET("/sync", runner.syncstatus)
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
	wallet.GET("/links", runner.links)
	wallet.POST("/links/request", runner.linkrequest)
	wallet.POST("/links/accept", runner.linkaccept)
	wallet.POST("/links/message", runner.sendlinkmessage)
	wallet.GET("/links/message", runner.getlinkmessages)
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
