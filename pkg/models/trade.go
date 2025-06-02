package models

import "time"


type Trade struct {
	ID            int64     `json:"id"`
	Symbol		string    `json:"symbol"`
	BuyOrderID  int64     `json:"buy_order_id"`
	SellOrderID int64     `json:"sell_order_id"`
	Price        float64   `json:"price"`
	Quantity     float64     `json:"quantity"`
	ExecutedAt   time.Time `json:"executed_at"`
}


