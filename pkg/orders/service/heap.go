package service

import (
	"container/heap"
	"fmt"
	"time"

	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/models"
)

// OrderHeap is a min-heap for orders, where the top of the heap is the order with the highest priority.
// For buy orders, the highest price is prioritized, and for sell orders, the lowest price is prioritized.

type Heap struct{
	orders []*models.Order
	less func(i, j *models.Order) bool
}


func (h *Heap) Len() int {
	return len(h.orders)
}


func (h *Heap) Less(i, j int) bool {
	return h.less(h.orders[i], h.orders[j])
}

func (h *Heap) Swap(i, j int) {
	h.orders[i], h.orders[j] = h.orders[j], h.orders[i]
}

func (h *Heap) Push(x interface{}) {
	h.orders = append(h.orders, x.(*models.Order))
}

func (h *Heap) Pop() interface{} {
	old := h.orders
	n := len(old)
	x := old[n-1]
	h.orders = old[:n-1]
	return x
}

func (h *Heap) Peek() interface{} {
	if h.Len() == 0 {
		return nil
	}
	return h.orders[0]
}

func (h *Heap) IsEmpty() bool {
	return len(h.orders) == 0
}


func (h *Heap) AddOrderToHeap(order *models.Order){
	if(order.Type == "limit"){
		heap.Push(h,order)
	}
}

func (h *Heap) RemoveTop() *models.Order {
	if h.IsEmpty() {
		return nil
	}
	order := heap.Pop(h).(*models.Order)
	return order
}

func (h *Heap) GetTopOrder() *models.Order {
	if h.IsEmpty() {
		return nil
	}
	order := h.Peek()
	return order.(*models.Order)
}

func (h *Heap) GetAllOrders() []*models.Order {
	return h.orders
}

func (h *Heap) UpdateTopOrder(updatedQuantity float64) *models.Order {
	if h.IsEmpty() {
		return nil
	}
	topOrder := h.Peek().(*models.Order)
	topOrder.RemainingQty = updatedQuantity 
	topOrder.Status = models.Partial
	topOrder.UpdatedAt = time.Now()
	if topOrder.RemainingQty <= 0 {
		h.RemoveTop()
	}else{
		heap.Fix(h, 0)
	}
	return topOrder
}

func BuildHeapFromData(array []*models.Order, side string) *Heap {
	orderheap := Heap{}
	if(side == "buy") {
		orderheap = Heap{
			orders: array,
			less: func(a, b *models.Order) bool {
				if a.Price == b.Price {
					return a.CreatedAt.Before(b.CreatedAt)
				}
				return a.Price > b.Price 
			},
		}
	}else if(side == "sell") {
		orderheap = Heap{
			orders: array,
			less: func(a, b *models.Order) bool {
				if a.Price == b.Price {
					return a.CreatedAt.Before(b.CreatedAt)
				}
				return a.Price < b.Price
			},
		}

	}
	heap.Init(&orderheap)
	fmt.Println("buildHeapByInit: ", orderheap.Peek())
	return &orderheap
}

func NewBuyHeap() *Heap {
	return &Heap{
		orders: []*models.Order{},
		less: func(a, b *models.Order) bool {
			if a.Price == b.Price {
				return a.CreatedAt.Before(b.CreatedAt)
			}
			return a.Price > b.Price 
		},
	}
}

func NewSellHeap() *Heap {
	return &Heap{
		orders: []*models.Order{},
		less: func(a, b *models.Order) bool {
			if a.Price == b.Price {
				return a.CreatedAt.Before(b.CreatedAt)
			}
			return a.Price < b.Price
		},
	}
}

func (h *Heap) PrintHeapState() {
	if h.IsEmpty() {
		fmt.Println("Heap is empty")
		return
	}
	fmt.Println("Current heap state:")
	for i, order := range h.orders {
		fmt.Printf("Index: %d, ID: %d, Price: %.2f, Qty: %d, CreatedAt: %s\n", 
			i, order.ID, order.Price, order.RemainingQty, order.CreatedAt.Format("15:04:05"))
	}
	fmt.Println()
}

func TestRunHeap(){
	buyHeap := NewBuyHeap()
	sellHeap := NewSellHeap()

	heap.Init(buyHeap)
	heap.Init(sellHeap)
	price95 := 95.0
    price100 := 100.0
    price105 := 105.0
    //price98 := 98.0
    //price102 := 102.0

	baseTime := time.Now()
    
    buyOrders := []*models.Order{
        {ID: 1, Side: "buy", Type: "limit", Price: price100, RemainingQty: (10), CreatedAt: baseTime.Add(2 * time.Second)},
        {ID: 2, Side: "buy", Type: "limit", Price: price105, RemainingQty: (8), CreatedAt: baseTime.Add(1 * time.Second)},
        {ID: 3, Side: "buy", Type: "limit", Price: price95, RemainingQty: (12), CreatedAt: baseTime},
        {ID: 4, Side: "buy", Type: "limit", Price: price100, RemainingQty: (6), CreatedAt: baseTime.Add(3 * time.Second)},
    }
	buyHeap = BuildHeapFromData(buyOrders, "buy")
	fmt.Println("Popping from Buy Heap:")
	
	for _, order := range buyOrders {
		sellHeap.AddOrderToHeap(order)
	}
	fmt.Println("sell Orders Heap:")
	sellHeap.PrintHeapState()
	for !sellHeap.IsEmpty() {
		order := sellHeap.RemoveTop()
		fmt.Printf("ID: %d, Price: %.2f, CreatedAt: %s\n", 
			order.ID, order.Price, order.CreatedAt.Format("15:04:05"))
	}
	fmt.Println("buy Orders Heap:")
	buyHeap.PrintHeapState()
	orderData := buyHeap.UpdateTopOrder(4)
	fmt.Println("All Orders in Buy Heap:", orderData)
	for !buyHeap.IsEmpty() {
		order := buyHeap.RemoveTop()
		fmt.Printf("ID: %d, Price: %.2f, CreatedAt: %s\n", 
			order.ID, order.Price, order.CreatedAt.Format("15:04:05"))
	}

	
}




