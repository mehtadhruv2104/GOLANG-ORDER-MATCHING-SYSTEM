package service

import (
	"database/sql"
	"fmt"

	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/models"
)

type OrderBookManager struct {
	OrderBooks map[string]*OrderBook
	DB         *sql.DB
}

type OrderBook struct {
	Symbol	  string
	BuyOrders  *Heap
	SellOrders *Heap
}

func NewOrderBookManager(db *sql.DB) *OrderBookManager {
	return &OrderBookManager{
		OrderBooks: make(map[string]*OrderBook),
		DB:         db,
	}
}

func (bookManager *OrderBookManager) GetOrCreateOrderBook(symbol string) *OrderBook {
	if orderBook, exists := bookManager.OrderBooks[symbol]; exists {
		return orderBook
	}
	
	orderBook := &OrderBook{
		Symbol:     symbol,
		BuyOrders:  NewBuyHeap(),
		SellOrders: NewSellHeap(),
	}
	bookManager.OrderBooks[symbol] = orderBook
	return orderBook
}

func (bookManager *OrderBookManager) SyncWithDatabase() error {
	fmt.Println("Database sync Start")
	query := `
        SELECT id, symbol, side, type, price, quantity, remaining_quantity, 
               status, created_at, updated_at
        FROM orders 
        WHERE status IN ('open', 'partial')
        ORDER BY symbol, side, 
                 CASE WHEN side = 'buy' THEN price END DESC,  -- Buy orders: highest price first
                 CASE WHEN side = 'sell' THEN price END ASC,  -- Sell orders: lowest price first
                 created_at ASC  -- Time priority for same price
    `
	rows, err := bookManager.DB.Query(query)
    if err != nil {
        return fmt.Errorf("failed to query orders: %w", err)
    }
	fmt.Println("rows",rows)
    defer rows.Close()
	if bookManager.OrderBooks == nil {
        bookManager.OrderBooks = make(map[string]*OrderBook)
    }
	orderCount := 0
    symbolCount := make(map[string]int)
	for rows.Next() {
        var order models.Order
        err := rows.Scan(
            &order.ID, &order.Symbol, &order.Side, &order.Type,
            &order.Price, &order.Quantity, &order.RemainingQty,
            &order.Status, &order.CreatedAt, &order.UpdatedAt,
        )
		fmt.Println("order", order, err)
		fmt.Println("book manager", bookManager.OrderBooks)
        if err != nil {
            return fmt.Errorf("failed to scan order: %w", err)
        }
        orderBook := bookManager.GetOrCreateOrderBook(order.Symbol)
		fmt.Println("orderBook", orderBook)
        orderBook.AddLimitOrder(order)
        orderCount++
        symbolCount[order.Symbol]++
    }
	fmt.Println("Orders Loaded to OrderBook", orderCount)
	fmt.Println("Symbols are there in the Books", symbolCount)
    
    if err = rows.Err(); err != nil {
        return fmt.Errorf("error iterating over rows: %w", err)
    }
	return nil
}

func (book *OrderBook) MatchMarketOrder(order *models.Order) (*models.Order, *models.Order, float64, error) {
	
	if order.Symbol != book.Symbol{
		return nil,order,0,fmt.Errorf("Wrong Symbol Matched")
	}
	order.RemainingQty = order.Quantity
	var matchedOrder *models.Order
	var quantity float64

	if order.Side == models.Buy {
		sellHeap := book.SellOrders
		topOrder := sellHeap.GetTopOrder()
		if topOrder == nil {
			return nil, order, 0, fmt.Errorf("no orders in sell order book")
		}
		if order.Type == models.Limit && topOrder.Price > order.Price {
			return nil, order, 0, fmt.Errorf("no orders matching in sell order book")
		}	
			
		if(order.Quantity == topOrder.RemainingQty){

			topOrder.RemainingQty = 0
			order.RemainingQty = 0
			quantity = order.Quantity
			topOrder.Status = models.Completed
			order.Status = models.Completed
			sellHeap.RemoveTop()

		}else if(order.Quantity > topOrder.RemainingQty){
			
			order.RemainingQty = order.Quantity - topOrder.RemainingQty
			quantity = topOrder.RemainingQty
			topOrder.RemainingQty = 0
			topOrder.Status = models.Completed
			order.Status = models.Partial
			sellHeap.RemoveTop()

		}else{
			
			topOrder.RemainingQty = topOrder.RemainingQty - order.Quantity
			order.RemainingQty = 0
			quantity = order.Quantity
			topOrder.Status = models.Partial
			order.Status = models.Completed
			sellHeap.UpdateTopOrder(topOrder.RemainingQty)

		}
		order.Price = topOrder.Price
		sellHeap.PrintHeapState()
		matchedOrder = topOrder

	}else if(order.Side == models.Sell){
		buyHeap := book.BuyOrders
		topOrder := buyHeap.GetTopOrder()
		if topOrder == nil {
			return nil,nil, 0, fmt.Errorf("no orders in buy order book")
		}
		if order.Type == models.Limit && topOrder.Price < order.Price{
			return nil, nil, 0, fmt.Errorf("no orders matching in buy order book")
		}

		if(order.Quantity == topOrder.RemainingQty){
			
			topOrder.RemainingQty = 0
			order.RemainingQty = 0
			quantity = order.Quantity
			topOrder.Status = models.Completed
			order.Status = models.Completed
			buyHeap.RemoveTop()

		}else if(order.Quantity > topOrder.RemainingQty){
			
			order.RemainingQty = order.Quantity - topOrder.RemainingQty
			quantity = topOrder.RemainingQty
			topOrder.RemainingQty = 0
			topOrder.Status = models.Completed
			order.Status = models.Partial
			buyHeap.RemoveTop()

		}else{
			
			topOrder.RemainingQty = topOrder.RemainingQty - order.Quantity
			order.RemainingQty = 0
			quantity = order.Quantity
			topOrder.Status = models.Partial
			order.Status = models.Completed
			buyHeap.UpdateTopOrder(topOrder.RemainingQty)
			
		}
		order.Price = topOrder.Price
		buyHeap.PrintHeapState()
		matchedOrder = topOrder	
	}
	return matchedOrder, order,quantity,nil
}

func (book *OrderBook) AddLimitOrder(order models.Order) (error) {
	if order.Symbol != book.Symbol{
		return fmt.Errorf("Wrong Symbol Matched")
	}
	if(order.Side == models.Buy){
		buyHeap := book.BuyOrders
		buyHeap.AddOrderToHeap(&order)
	}else if(order.Side == models.Sell){
		sellHeap := book.SellOrders
		sellHeap.AddOrderToHeap(&order)
	}
	return nil
}







