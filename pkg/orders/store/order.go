package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM/pkg/models"
)

type OrderStore interface {
	AddOrder(order *models.Order) (int64,error)
	GetOrderByID(id int64) (*models.Order, error)
	GetOrdersBySymbol(symbol string) ([]models.Order, error) 
	GetOrdersByStatus(symbol string, status string) ([]models.Order, error)
	RemoveOrder(id int64) error
	UpdateOrder(order *models.Order) error
	CancelOrder(id int64) error
}


type OrderStoreManager struct {
	DB *sql.DB
}

func NewOrderStore(db *sql.DB) OrderStore {
	return &OrderStoreManager{
		DB: db,
	}
}

func (store *OrderStoreManager) AddOrder(order *models.Order) (int64,error) {

	query := `INSERT INTO orders (symbol, side, type, price, quantity, remaining_quantity, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	order.RemainingQty = order.Quantity 
	if order.Status == "" {
		order.Status = models.Open
	}
	result,err := store.DB.Exec(
		query,
		order.Symbol,
		order.Side,
		order.Type,
		order.Price,
		order.Quantity,
		order.RemainingQty,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		fmt.Println("Error adding order to database:", err)
		return 0,err
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Failed to fetch last insert ID:", err)
		return id,err
	}
	order.ID = id 
	return id,nil
}


func (store *OrderStoreManager) GetOrderByID(id int64) (*models.Order, error) {
	query := `
		SELECT id, symbol, side, type, price, quantity, remaining_quantity, status, created_at, updated_at
		FROM orders 
		WHERE id = ?
	`
	
	order := &models.Order{}
	err := store.DB.QueryRow(query, id).Scan(
		&order.ID,
		&order.Symbol,
		&order.Side,
		&order.Type,
		&order.Price,
		&order.Quantity,
		&order.RemainingQty,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	return order, nil
}

func (store *OrderStoreManager) GetOrdersBySymbol(symbol string) ([]models.Order, error) {
	query := `
		SELECT id, symbol, side, type, price, quantity, remaining_quantity, status, created_at, updated_at
		FROM orders 
		WHERE symbol = ?
		ORDER BY created_at ASC
	`
	return store.queryOrders(query, symbol)
}

func (store *OrderStoreManager) GetOrdersByStatus(symbol string, status string) ([]models.Order, error) {
	query := `
		SELECT id, symbol, side, type, price, quantity, remaining_quantity, status, created_at, updated_at
		FROM orders 
		WHERE symbol = ? AND status = ?
		ORDER BY created_at ASC
	`
	
	return store.queryOrders(query, symbol, status)
}

func (store *OrderStoreManager) UpdateOrder(order *models.Order) error {
	query := `
		UPDATE orders 
		SET remaining_quantity = ?, status = ?, updated_at = ?
		WHERE id = ?
	`
	
	order.UpdatedAt = time.Now()
	
	result, err := store.DB.Exec(query, order.RemainingQty, order.Status, order.UpdatedAt, order.ID)
	if err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found or no changes made")
	}

	return nil
}

func (osm *OrderStoreManager) CancelOrder(id int64) error {
	query := `
		UPDATE orders 
		SET status = 'canceled', updated_at = ?
		WHERE id = ? AND status IN ('open', 'partial')
	`
	
	result, err := osm.DB.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found or cannot be canceled (already filled/canceled)")
	}

	return nil
}

func (store *OrderStoreManager) RemoveOrder(id int64) error {
	query := `DELETE FROM orders WHERE id = ?`
	
	result, err := store.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to remove order: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

func (store *OrderStoreManager) queryOrders(query string, args ...interface{}) ([]models.Order, error) {
	rows, err := store.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	return store.scanOrders(rows)
}

func (store *OrderStoreManager) scanOrders(rows *sql.Rows) ([]models.Order, error) {
	var orders []models.Order

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.Symbol,
			&order.Side,
			&order.Type,
			&order.Price,
			&order.Quantity,
			&order.RemainingQty,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return orders, nil
}