package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
)



func (h OrderHandler) GetOrderBook(c *gin.Context) {
	symbol := c.Query("symbol")
	if strings.TrimSpace(symbol) == "" {
		c.JSON(400, gin.H{
			"error":"Symbol Not given",
		})
		return
	}
	orderBookResponse,err := h.OrderService.OrderBookManager.GetOrderBook(symbol)
	if err!=nil{
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	
	c.JSON(200,orderBookResponse)

}