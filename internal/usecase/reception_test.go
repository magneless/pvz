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

type mockReceptionRepo struct {
	mock.Mock
}

func (m *mockReceptionRepo) Create(ctx context.Context, reception models.Reception) error {
	return m.Called(ctx, reception).Error(0)
}

func (m *mockReceptionRepo) CloseReception(ctx context.Context, pvzID uuid.UUID) error {
	return m.Called(ctx, pvzID).Error(0)
}

func TestReceptionUsecase_CreateReception(t *testing.T) {
	repo := new(mockReceptionRepo)
	uc := usecase.NewReceptionUsecase(repo)

	pvzID := uuid.New()
	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	reception, err := uc.CreateReception(context.Background(), pvzID)
	require.NoError(t, err)
	require.Equal(t, pvzID, reception.PVZID)
}

func TestReceptionUsecase_CloseReception(t *testing.T) {
	repo := new(mockReceptionRepo)
	uc := usecase.NewReceptionUsecase(repo)

	pvzID := uuid.New()
	repo.On("CloseReception", mock.Anything, pvzID).Return(nil)

	err := uc.CloseReception(context.Background(), pvzID)
	require.NoError(t, err)
}
