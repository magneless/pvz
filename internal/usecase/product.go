package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/domain/repository"
)

var ErrWrongType = errors.New("wrong type")

type ProductUsecase struct {
	repository repository.ProductRepository
}

func NewProductUsecase(repository repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		repository: repository,
	}
}

func (uc *ProductUsecase) CreateProduct(ctx context.Context, pType models.Type, pvzID uuid.UUID) (models.Product, error) {
	if !pType.IsValid() {
		return models.Product{}, ErrWrongType
	}
	now := time.Now()
	product := models.Product{
		ID:       uuid.New(),
		DateTime: &now,
		Type:     pType,
	}
	receptionID, err := uc.repository.Create(ctx, product, pvzID)
	if err != nil {
		return models.Product{}, err
	}
	product.ReceptionID = receptionID

	return product, nil
}

func (uc *ProductUsecase) DeleteProduct(ctx context.Context, pvzID uuid.UUID) error {
	return uc.repository.DeleteProduct(ctx, pvzID)
}