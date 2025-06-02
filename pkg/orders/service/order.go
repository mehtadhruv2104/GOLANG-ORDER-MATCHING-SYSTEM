package service

import (
	"errors"
	"fmt"

	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/models"
	orderstore "github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/orders/store"
	tradestore "github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/trades/store"
)





type OrderService interface{
	PlaceOrder(order models.Order) (models.Order, error)
	AddOrder(order *models.Order) (int64,error)
	AddTrade(orderA *models.Order, orderB *models.Order, quantity float64) (error)
}


type Service struct {
	OrderBookManager *OrderBookManager
	OrderStore orderstore.OrderStore
	TradeStore tradestore.TradeStore
}


func NewOrderService(store orderstore.OrderStore, orderBookManager *OrderBookManager, tradestore tradestore.TradeStore) *Service {
	return &Service{
		OrderBookManager: orderBookManager,
		OrderStore: store,
		TradeStore: tradestore,
	}
}

func (s *Service) AddOrder(order *models.Order) (int64,error) {
	
	order.RemainingQty = order.Quantity
	id,err := s.OrderStore.AddOrder(order)
	if err != nil {
		return 0,errors.New("failed to add order"+ err.Error())
	}
	return id,nil
}

func (s *Service) AddTrade(orderA *models.Order, orderB *models.Order, quantity float64) (error) {
	var buyOrder,sellOrder *models.Order
	if(orderA.Side == models.Buy){
		buyOrder = orderA
		sellOrder = orderB
	}else{
		buyOrder = orderB
		sellOrder = orderA
	}

	trade := models.Trade{
		BuyOrderID: buyOrder.ID,
		SellOrderID: sellOrder.ID,
		Symbol: sellOrder.Symbol,
		Price: sellOrder.Price,
		Quantity: quantity,
	}

	id,err := s.TradeStore.AddTrade(&trade)
	if err != nil {
		return errors.New("failed to execute trade"+ err.Error())
	}
	trade.ID = id
	return nil
}

func (s *Service) PlaceOrder(order models.Order) (error) {
	
	orderBook := s.OrderBookManager.GetOrCreateOrderBook(order.Symbol)

	if( order.Type == models.Market) {

		matchedOrder, updatedOrder, quantity, err := orderBook.MatchMarketOrder(&order)
		if err != nil {
			return errors.New("Market order could not be completed."+ err.Error())
		}

		id,err := s.AddOrder(updatedOrder)
		if err != nil {
			return errors.New("Failed to add order to database"+err.Error())
		}
		order.ID = id

		err = s.AddTrade(updatedOrder,matchedOrder, quantity)
		if err != nil {
			return errors.New("Trade could not be executed" + err.Error())
		}

		err = s.OrderStore.UpdateOrder(matchedOrder)
		if err != nil {
			return errors.New("Failed to update order to database"+err.Error())
		}


	} else if(order.Type == models.Limit) {

		id,err := s.AddOrder(&order)
		if err != nil {
			return errors.New("Failed to add order to database"+ err.Error())
		}
		order.ID = id
		
		matchedOrder, updatedOrder, quantity, err := orderBook.MatchMarketOrder(&order)
		if err == nil {
			err = s.AddTrade(updatedOrder,matchedOrder, quantity)
			if err != nil {
				return errors.New("Trade could not be executed"+err.Error())
			}

			err = s.OrderStore.UpdateOrder(matchedOrder)
			if err != nil {
				return errors.New("Failed to update order to database"+ err.Error())
			}

			err = s.OrderStore.UpdateOrder(updatedOrder)
			if err != nil {
				return errors.New("Failed to update order to database"+ err.Error())
			}
		}

		if order.Status != models.Completed{
			err = orderBook.AddLimitOrder(order)
			if err != nil {
				return errors.New("Failed to add limit order" + err.Error())
			}
		}
		
	}
	return nil
}

func (s *Service) PrintCurrentHeapState(symbol string){
	orderBook := s.OrderBookManager.GetOrCreateOrderBook(symbol)
	fmt.Println("Buy Heap Top", orderBook.BuyOrders.GetTopOrder())
	fmt.Println("Sell Heap Top", orderBook.SellOrders.GetTopOrder())
}

func (s *Service) GetOrderByID(orderID int64)(*models.Order, error){
	var orderDemo models.Order
	if orderID == 0{
		return &orderDemo,errors.New("invalid Order Id")
	}
	order,err :=s.OrderStore.GetOrderByID(orderID)
	if err != nil{
		return order,errors.New(err.Error())
	}
	return order,nil
}

func (s *Service) CancelOrder(orderID int64)(error){
	if orderID == 0{
		return errors.New("invalid Order Id")
	}
	err :=s.OrderStore.CancelOrder(orderID)
	if err != nil{
		return errors.New("order Referencing this orderId not found")
	}
	return nil
}



