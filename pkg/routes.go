package pkg

import (
	"github.com/gin-gonic/gin"
	orderHandler "github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/orders/handler"
)




func InitTradeEngineRoutes(apiRoutes *gin.RouterGroup, h *orderHandler.OrderHandler) {

	apiRoutes.POST("/orders", h.PlaceOrder)
	apiRoutes.DELETE("/orders/:id", h.CancelOrder)
	apiRoutes.GET("/orders/:id", h.GetOrderStatus)
	apiRoutes.GET(("/orderbook"), h.GetOrderBook)
	apiRoutes.GET("/trades", h.GetTrades)
}