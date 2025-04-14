package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
)

type ProductRepository interface {
	Create(ctx context.Context, product models.Product, pvzID uuid.UUID) (uuid.UUID, error)
	DeleteProduct(ctx context.Context, pvzID uuid.UUID) error
}