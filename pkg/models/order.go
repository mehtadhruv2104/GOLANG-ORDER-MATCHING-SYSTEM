package models

import "time"

type Order struct{
	ID          	int64    		`json:"id"`
	Symbol 			string  		`json:"symbol"`
	Side			OrderSide  		`json:"side"`
	Type 			OrderType  		`json:"type"`
	Price 			float64  		`json:"price"`
	Quantity 		int64  			`json:"quantity"`
	RemainingQty 	int64  			`json:"remaining_qty"`
	Status 			OrderStatus  	`json:"status"`
	CreatedAt 		time.Time  		`json:"created_at"`
	UpdatedAt 		time.Time  		`json:"updated_at"`
}

type OrderSide string
const (
	Buy  OrderSide = "buy"
	Sell OrderSide = "sell"
)

type OrderType string
const (
	Market OrderType = "market"
	Limit  OrderType = "limit"
)

type OrderStatus string
const (
	Open   OrderStatus = "open"
	Completed OrderStatus = "completed"
	Cancelled OrderStatus = "cancelled"
	Partial OrderStatus = "partial"
)