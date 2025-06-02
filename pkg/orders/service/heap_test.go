package service_test

import (
	"container/heap"
	"testing"
	"time"

	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/models"
	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/orders/service"
	"github.com/stretchr/testify/assert"
)

func TestBuyAndSellHeaps(t *testing.T) {
	price95 := 95.0
	price100 := 100.0
	price105 := 105.0

	baseTime := time.Now()

	buyOrders := []*models.Order{
		{ID: 1, Side: "buy", Type: "limit", Price: price100, RemainingQty: 10, CreatedAt: baseTime.Add(2 * time.Second)},
		{ID: 2, Side: "buy", Type: "limit", Price: price105, RemainingQty: 8, CreatedAt: baseTime.Add(1 * time.Second)},
		{ID: 3, Side: "buy", Type: "limit", Price: price95, RemainingQty: 12, CreatedAt: baseTime},
		{ID: 4, Side: "buy", Type: "limit", Price: price100, RemainingQty: 6, CreatedAt: baseTime.Add(3 * time.Second)},
	}

	expectedBuyOrderIDs := []int64{2, 1, 4, 3}
	expectedSellOrderIDs := []int64{3, 1, 4, 2}

	// --- Test Buy Heap (Max-Heap) ---
	buyHeap := service.NewBuyHeap()
	heap.Init(buyHeap)

	for _, order := range buyOrders {
		buyHeap.AddOrderToHeap(order)
	}

	var actualBuyOrderIDs []int64
	for !buyHeap.IsEmpty() {
		order := buyHeap.RemoveTop()
		actualBuyOrderIDs = append(actualBuyOrderIDs, order.ID)
	}

	assert.Equal(t, expectedBuyOrderIDs, actualBuyOrderIDs, "Buy heap order mismatch")

	// --- Test Sell Heap (Min-Heap) ---
	sellHeap := service.NewSellHeap()
	heap.Init(sellHeap)

	for _, order := range buyOrders {
		sellHeap.AddOrderToHeap(order)
	}

	var actualSellOrderIDs []int64
	for !sellHeap.IsEmpty() {
		order := sellHeap.RemoveTop()
		actualSellOrderIDs = append(actualSellOrderIDs, order.ID)
	}

	assert.Equal(t, expectedSellOrderIDs, actualSellOrderIDs, "Sell heap order mismatch")
}
