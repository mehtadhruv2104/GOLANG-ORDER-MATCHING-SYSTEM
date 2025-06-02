package handler

import (
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/orders/service"
)




type OrderHandler struct {
	OrderService *service.Service
}

func NewOrderhandler(orderService *service.Service) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}




