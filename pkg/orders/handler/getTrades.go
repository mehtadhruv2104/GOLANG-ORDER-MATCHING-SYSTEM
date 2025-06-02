package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
)


func (h OrderHandler) GetTrades(c *gin.Context) {
	symbol := c.Query("symbol")
	if strings.TrimSpace(symbol) == "" {
		c.JSON(400, gin.H{
			"error":"Symbol Not given",
		})
		return
	}
	trades,err := h.OrderService.TradeStore.GetTradesBySymbol(symbol)
	if err != nil{
		c.JSON(500, gin.H{
			"error":"Could not find trades",
		})
	}
	c.JSON(200,trades)

}