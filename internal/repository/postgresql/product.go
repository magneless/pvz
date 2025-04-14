package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/domain/repository"
)

var ErrNotOpenReception = errors.New("not open reception")


type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product models.Product, pvzID uuid.UUID) (uuid.UUID, error) {
	const query = `
		INSERT INTO products (id, date_time, type, reception_id)
		VALUES (
			$1,
			$2,
			$3, 
			(SELECT id FROM receptions WHERE pvz_id = $4 AND status = 'in_progress' LIMIT 1)
		)
		RETURNING reception_id;
	`
	var receptionID uuid.UUID

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.ID,
		*product.DateTime,
		product.Type,
		pvzID,
	).Scan(&receptionID)

	if err == sql.ErrNoRows {
		return uuid.Nil, ErrNotOpenReception
	}
	if err != nil {
		return uuid.Nil, err
	}

	return receptionID, err
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, pvzID uuid.UUID) error {
	const query = `
    WITH active_reception AS (
        SELECT id AS reception_id
        FROM receptions
        WHERE pvz_id = $1
          AND status = 'in_progress'
        LIMIT 1
    ),
    last_product AS (
        SELECT id AS product_id
        FROM products
        WHERE reception_id = (SELECT reception_id FROM active_reception)
        ORDER BY date_time DESC
        LIMIT 1
    )
    DELETE FROM products
    WHERE id = (SELECT product_id FROM last_product)
    RETURNING id;
    `
	var deletedID uuid.UUID
    err := r.db.QueryRowContext(ctx, query, pvzID).Scan(&deletedID)
    return err
}
