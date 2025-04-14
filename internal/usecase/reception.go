package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/domain/repository"
)

type ReceptionUsecase struct {
	repository repository.ReceptionRepository
}

func NewReceptionUsecase(repository repository.ReceptionRepository) *ReceptionUsecase {
	return &ReceptionUsecase{
		repository: repository,
	}
}

func (uc *ReceptionUsecase) CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	now := time.Now()
	reception := models.Reception{
		ID:       uuid.New(),
		DateTime: &now,
		PVZID:    pvzID,
		Status:   "in_progress",
	}
	err := uc.repository.Create(ctx, reception)
	return reception, err
}

func (uc *ReceptionUsecase) CloseReception(ctx context.Context, pvzID uuid.UUID) error {
	return uc.repository.CloseReception(ctx, pvzID)
}