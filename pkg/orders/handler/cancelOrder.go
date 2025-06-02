package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)




func (h OrderHandler) CancelOrder(c *gin.Context) {

	orderID := c.Param("id")
	ID,err := strconv.ParseInt(orderID,10,64)
	if err!= nil{
		c.JSON(400,gin.H{"error": "Failed to convert ID string to int" + err.Error()})
		return
	}
	err = h.OrderService.CancelOrder(ID)
	if err != nil{
		c.JSON(400,gin.H{"error": "Could not find the order" + err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "The Order has been cancelled successfully",
	})


}
