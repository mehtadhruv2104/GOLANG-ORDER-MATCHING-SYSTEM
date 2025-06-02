package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/models"
)


type TradeStore interface {
	AddTrade(trade *models.Trade) (int64,error)
	GetTradesBySymbol(symbol string) ([]models.Trade, error)
	GetTradeByID(id int64) (*models.Trade, error)
	GetTradesByOrderID(orderID int64) ([]models.Trade, error) 
}

type TradeStoreManager struct {
	DB *sql.DB
}

func NewTradeStore(db *sql.DB) *TradeStoreManager {
	return &TradeStoreManager{
		DB: db,
	}
}

func (store *TradeStoreManager) AddTrade(trade *models.Trade) (int64,error) {
	query := `
		INSERT INTO trades (symbol, buy_order_id, sell_order_id, price, quantity, executed_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	
	if trade.ExecutedAt.IsZero() {
		trade.ExecutedAt = time.Now()
	}

	result,err := store.DB.Exec(
		query,
		trade.Symbol,
		trade.BuyOrderID,
		trade.SellOrderID,
		trade.Price,
		trade.Quantity,
		trade.ExecutedAt,
	)
	if err != nil {
		return 0,fmt.Errorf("failed to add trade: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Failed to fetch last insert ID:", err)
		return 0,err
	}
	trade.ID = id
	return id,nil
}


func (store *TradeStoreManager) GetTradeByID(id int64) (*models.Trade, error) {
	query := `
		SELECT id, symbol, buy_order_id, sell_order_id, price, quantity, executed_at
		FROM trades 
		WHERE id = ?
	`
	
	trade := &models.Trade{}
	err := store.DB.QueryRow(query, id).Scan(
		&trade.ID,
		&trade.Symbol,
		&trade.BuyOrderID,
		&trade.SellOrderID,
		&trade.Price,
		&trade.Quantity,
		&trade.ExecutedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade not found")
		}
		return nil, fmt.Errorf("failed to get trade: %v", err)
	}

	return trade, nil
}

func (store *TradeStoreManager) GetTradesBySymbol(symbol string) ([]models.Trade, error) {
	query := `
		SELECT id, symbol, buy_order_id, sell_order_id, price, quantity, executed_at
		FROM trades 
		WHERE symbol = ?
		ORDER BY executed_at DESC
	`
	
	return store.queryTrades(query, symbol)
}

func (store *TradeStoreManager) GetRecentTradesBySymbol(symbol string, limit int) ([]models.Trade, error) {
	query := `
		SELECT id, symbol, buy_order_id, sell_order_id, price, quantity, executed_at
		FROM trades 
		WHERE symbol = ?
		ORDER BY executed_at DESC
		LIMIT ?
	`
	return store.queryTrades(query, symbol, limit)
}

func (store *TradeStoreManager) GetTradesByOrderID(orderID int64) ([]models.Trade, error) {
	query := `
		SELECT id, symbol, buy_order_id, sell_order_id, price, quantity, executed_at
		FROM trades 
		WHERE buy_order_id = ? OR sell_order_id = ?
		ORDER BY executed_at DESC
	`
	
	return store.queryTrades(query, orderID)
}

func (store *TradeStoreManager) queryTrades(query string, args ...interface{}) ([]models.Trade, error) {
	rows, err := store.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %v", err)
	}
	defer rows.Close()

	return store.scanTrades(rows)
}

func (store *TradeStoreManager) scanTrades(rows *sql.Rows) ([]models.Trade, error) {
	var trades []models.Trade

	for rows.Next() {
		var trade models.Trade
		err := rows.Scan(
			&trade.ID,
			&trade.Symbol,
			&trade.BuyOrderID,
			&trade.SellOrderID,
			&trade.Price,
			&trade.Quantity,
			&trade.ExecutedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %v", err)
		}
		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return trades, nil
}

