package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/domain/repository"
)

type PVZUsecase struct {
	repository repository.PVZRepository
}

func NewPVZUsecase(repository repository.PVZRepository) *PVZUsecase {
	return &PVZUsecase{
		repository: repository,
	}
}

func (uc *PVZUsecase) CreatePVZ(ctx context.Context, pvz models.PVZ) error {
	if !pvz.City.IsValid() {
		return errors.New("wrong city")
	}
	if pvz.ID == uuid.Nil {
		pvz.ID = uuid.New()
	}
	if pvz.RegistrationDate == nil {
		now := time.Now()
		pvz.RegistrationDate = &now
	}
	return uc.repository.Create(ctx, pvz)
}

func (uc *PVZUsecase) GetPVZFullInfo(ctx context.Context, startDate, endDate *time.Time, page, limit *int) ([]models.PVZFullInfo, error) {
	sd := time.Time{}
	ed := time.Now()
	p := 1
	l := 10
	
	if startDate != nil {
		sd = *startDate
	}
	if endDate != nil {
		ed = *endDate
	}
	if ed.Compare(sd) == -1{
		return nil, errors.New("end date is more than start date")
	}

	if page != nil {
		p = *page
	}
	if p < 1 {
		return nil, errors.New("amount of pages less than 1")
	}

	if limit != nil {
		l = *limit
	}
	if l < 1 || l > 30 {
		return nil, errors.New("limit less than 1 or more than 30")
	}

	return uc.repository.GetAllInfo(ctx, sd, ed, p, l)
}