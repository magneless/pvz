package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
)

type ReceptionRepository interface {
	Create(ctx context.Context, reception models.Reception) error
	CloseReception(ctx context.Context, pvzID uuid.UUID) error
}
