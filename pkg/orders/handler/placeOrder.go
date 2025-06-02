package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/models"
)



func (h OrderHandler) PlaceOrder(c *gin.Context) {
	req := models.Order{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	err := h.OrderService.PlaceOrder(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to place order" + err.Error()})
		return
	}
	c.JSON(200, gin.H{"order": "Order placed successfully"})

}





