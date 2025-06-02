package pkg

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/orders/handler"
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/orders/service"
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/orders/store"
	tradestore "github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/trades/store"
)



func StartEngine(dB *sql.DB) (*gin.Engine) {
	router := gin.Default()
	router.Use(gin.Logger())
	apiRoutes := router.Group("/api")
	orderBookManager := service.NewOrderBookManager(dB)
	orderBookManager.SyncWithDatabase()
	orderStore := store.NewOrderStore(dB)
	tradeStore := tradestore.NewTradeStore(dB)
	orderService := service.NewOrderService(orderStore, orderBookManager, tradeStore)
	orderHandler := handler.NewOrderhandler(orderService)
	InitTradeEngineRoutes(apiRoutes, orderHandler)
	//apiRoutes.Use(middleware)
	//service.TestRunHeap()
	
	return router
}

