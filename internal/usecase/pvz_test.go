package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPVZRepo struct {
	mock.Mock
}

func (m *mockPVZRepo) Create(ctx context.Context, pvz models.PVZ) error {
	return m.Called(ctx, pvz).Error(0)
}

func (m *mockPVZRepo) GetAllInfo(ctx context.Context, sd, ed time.Time, page, limit int) ([]models.PVZFullInfo, error) {
	args := m.Called(ctx, sd, ed, page, limit)
	return args.Get(0).([]models.PVZFullInfo), args.Error(1)
}

func TestPVZUsecase_CreatePVZ(t *testing.T) {
	repo := new(mockPVZRepo)
	uc := usecase.NewPVZUsecase(repo)

	pvz := models.PVZ{City: "Москва"}

	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := uc.CreatePVZ(context.Background(), pvz)
	require.NoError(t, err)
}

func TestPVZUsecase_GetPVZFullInfo(t *testing.T) {
	repo := new(mockPVZRepo)
	uc := usecase.NewPVZUsecase(repo)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	expected := []models.PVZFullInfo{{}}
	repo.On("GetAllInfo", mock.Anything, mock.Anything, mock.Anything, 1, 10).Return(expected, nil)

	info, err := uc.GetPVZFullInfo(context.Background(), &start, &end, nil, nil)
	require.NoError(t, err)
	require.Len(t, info, 1)
}
