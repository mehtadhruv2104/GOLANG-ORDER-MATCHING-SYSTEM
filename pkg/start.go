package pkg

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)





func StartEngine(dB *sql.DB) (*gin.Engine) {
	router := gin.Default()
	router.Use(gin.Logger())
	apiRoutes := router.Group("/api")
	InitTradeEngineRoutes(apiRoutes)
	//apiRoutes.Use(middleware)
	return router
}