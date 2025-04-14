package repository

import (
	"context"
	"time"

	"github.com/magneless/pvz/internal/domain/models"
)

type PVZRepository interface {
	Create(ctx context.Context, pvz models.PVZ) error
	GetAllInfo(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]models.PVZFullInfo, error)
}
