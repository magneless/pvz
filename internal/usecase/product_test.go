package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Мок репозитория
type mockProductRepo struct {
	mock.Mock
}

func (m *mockProductRepo) Create(ctx context.Context, product models.Product, pvzID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, product, pvzID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockProductRepo) DeleteProduct(ctx context.Context, pvzID uuid.UUID) error {
	return m.Called(ctx, pvzID).Error(0)
}

func TestProductUsecase_CreateProduct(t *testing.T) {
	repo := new(mockProductRepo)
	uc := usecase.NewProductUsecase(repo)

	ctx := context.Background()
	pvzID := uuid.New()
	pType := models.Type("электроника")
	receptionID := uuid.New()

	repo.On("Create", mock.Anything, mock.Anything, pvzID).Return(receptionID, nil)

	product, err := uc.CreateProduct(ctx, pType, pvzID)
	require.NoError(t, err)
	require.Equal(t, receptionID, product.ReceptionID)
	require.Equal(t, pType, product.Type)
}

func TestProductUsecase_CreateProduct_InvalidType(t *testing.T) {
	uc := usecase.NewProductUsecase(nil)
	_, err := uc.CreateProduct(context.Background(), models.Type("wrong"), uuid.New())
	require.ErrorIs(t, err, usecase.ErrWrongType)
}

func TestProductUsecase_DeleteProduct(t *testing.T) {
	repo := new(mockProductRepo)
	uc := usecase.NewProductUsecase(repo)

	pvzID := uuid.New()
	repo.On("DeleteProduct", mock.Anything, pvzID).Return(nil)

	err := uc.DeleteProduct(context.Background(), pvzID)
	require.NoError(t, err)
}
