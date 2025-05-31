package pkg

import (
	"github.com/gin-gonic/gin"
	orderHandler "github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/orders/handler"
	tradeHandler "github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/trades/handler"
)




func InitTradeEngineRoutes(apiRoutes *gin.RouterGroup) {

	apiRoutes.POST("/orders", orderHandler.PlaceOrder)
	apiRoutes.DELETE("/orders/:id", orderHandler.CancelOrder)
	apiRoutes.GET("/orders/:id", orderHandler.GetOrderStatus)
	apiRoutes.GET(("/orderbook"), orderHandler.GetOrderBook)
	apiRoutes.GET("/trades", tradeHandler.GetTrades)
}