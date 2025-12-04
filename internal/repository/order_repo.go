package repository

import (
	"database/sql"
	"encoding/json"
	"ta/internal/domain"
	"time"

	"github.com/google/uuid"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (id, user_id) VALUES ($1, $2)`
	_, err = tx.Exec(query, order.ID, order.UserID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		itemQuery := `INSERT INTO order_items (id, order_id, product_id, quantity, price_at_time) 
		              VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(itemQuery, item.ID, item.OrderID, item.ProductID, item.Quantity, item.PriceAtTime)
		if err != nil {
			return err
		}
	}

	snapshot, err := json.Marshal(order)
	if err != nil {
		return err
	}

	historyQuery := `INSERT INTO order_history (order_id, user_id, snapshot_data) VALUES ($1, $2, $3)`
	_, err = tx.Exec(historyQuery, order.ID, order.UserID, snapshot)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*domain.Order, error) {
	query := `SELECT id, user_id, created_at FROM orders WHERE id = $1`
	var orderID, userID uuid.UUID
	var createdAt time.Time
	err := r.db.QueryRow(query, id).Scan(&orderID, &userID, &createdAt)
	if err != nil {
		return nil, err
	}

	itemsQuery := `SELECT id, order_id, product_id, quantity, price_at_time, created_at 
	               FROM order_items WHERE order_id = $1`
	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.PriceAtTime, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &domain.Order{
		ID:        orderID,
		UserID:    userID,
		Items:     items,
		CreatedAt: createdAt,
	}, nil
}

func (r *OrderRepository) GetByUserID(userID uuid.UUID) ([]*domain.Order, error) {
	query := `SELECT id, user_id, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var orderID, uid uuid.UUID
		var createdAt time.Time
		if err := rows.Scan(&orderID, &uid, &createdAt); err != nil {
			return nil, err
		}

		itemsQuery := `SELECT id, order_id, product_id, quantity, price_at_time, created_at 
		               FROM order_items WHERE order_id = $1`
		itemRows, err := r.db.Query(itemsQuery, orderID)
		if err != nil {
			return nil, err
		}

		var items []domain.OrderItem
		for itemRows.Next() {
			var item domain.OrderItem
			if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.PriceAtTime, &item.CreatedAt); err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, item)
		}
		itemRows.Close()

		orders = append(orders, &domain.Order{
			ID:        orderID,
			UserID:    uid,
			Items:     items,
			CreatedAt: createdAt,
		})
	}

	return orders, rows.Err()
}

