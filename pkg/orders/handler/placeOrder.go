package handler

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/models"
)

func ValidateRequest(order models.Order) error{
	if order.Type == models.Limit && order.Price <= 0{
		return errors.New("price must be greater than 0 for limit orders")
	}
	if order.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if order.Side != models.Buy && order.Side != models.Sell {
		return errors.New("invalid side: must be 'buy' or 'sell'")
	}
	if order.Type != models.Market && order.Type != models.Limit {
		return errors.New("invalid type: must be 'market' or 'limit'")
	}
	if strings.TrimSpace(order.Symbol) == "" {
		return errors.New("symbol is required")
	}
	return nil
}

func (h OrderHandler) PlaceOrder(c *gin.Context) {
	req := models.Order{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	err := ValidateRequest(req)
	if err != nil{
		c.JSON(400, gin.H{"error": "Bad Request" + err.Error()})
		return
	}
	err = h.OrderService.PlaceOrder(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to place order" + err.Error()})
		return
	}
	c.JSON(200, gin.H{"order": "Order placed successfully"})

}





