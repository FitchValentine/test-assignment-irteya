package repository

import (
	"database/sql"
	"ta/internal/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

type ProductDB struct {
	ID          uuid.UUID
	Description string
	Tags        pq.StringArray
	Quantity    int
}

func (r *ProductRepository) Create(product *domain.Product) error {
	query := `INSERT INTO products (id, description, tags, quantity) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, product.ID, product.Description, pq.Array(product.Tags), product.Quantity)
	return err
}

func (r *ProductRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	query := `SELECT id, description, tags, quantity FROM products WHERE id = $1`
	var dbProduct ProductDB
	err := r.db.QueryRow(query, id).Scan(
		&dbProduct.ID, &dbProduct.Description, &dbProduct.Tags, &dbProduct.Quantity,
	)
	if err != nil {
		return nil, err
	}
	return toProductDomain(&dbProduct), nil
}

func (r *ProductRepository) UpdateQuantity(id uuid.UUID, quantity int) error {
	query := `UPDATE products SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, quantity, id)
	return err
}

func (r *ProductRepository) GetByIDs(ids []uuid.UUID) ([]*domain.Product, error) {
	if len(ids) == 0 {
		return []*domain.Product{}, nil
	}
	query := `SELECT id, description, tags, quantity FROM products WHERE id = ANY($1)`
	rows, err := r.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var dbProduct ProductDB
		if err := rows.Scan(&dbProduct.ID, &dbProduct.Description, &dbProduct.Tags, &dbProduct.Quantity); err != nil {
			return nil, err
		}
		products = append(products, toProductDomain(&dbProduct))
	}
	return products, rows.Err()
}

func toProductDomain(db *ProductDB) *domain.Product {
	return &domain.Product{
		ID:          db.ID,
		Description: db.Description,
		Tags:        []string(db.Tags),
		Quantity:    db.Quantity,
	}
}

