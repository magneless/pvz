package postgresql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/domain/repository"
)

type ReceptionRepository struct {
	db *sql.DB
}

func NewReceptionRepository(db *sql.DB) repository.ReceptionRepository {
	return &ReceptionRepository{
		db: db,
	}
}

func (r *ReceptionRepository) Create(ctx context.Context, reception models.Reception) error {
	const query = `
		INSERT INTO receptions (id, date_time, status, pvz_id)
		SELECT $1, $2, $3, $4
		WHERE NOT EXISTS (
			SELECT 1
			FROM receptions
			WHERE pvz_id = $4
			AND status = 'in_progress'
		)
		RETURNING id;
	`
	var insertedID uuid.UUID
	err := r.db.QueryRowContext(
		ctx,
		query,
		reception.ID,
		*reception.DateTime,
		reception.Status,
		reception.PVZID,
	).Scan(&insertedID)
	
	return err
}

func (r *ReceptionRepository) CloseReception(ctx context.Context, pvzID uuid.UUID) error {
	const query = `
		UPDATE receptions
		SET status = 'closed'
		WHERE id = (
			SELECT id FROM receptions
			WHERE pvz_id = $1 AND status <> 'closed'
			ORDER BY date_time DESC
			LIMIT 1
		)
		RETURNING id;
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, pvzID).Scan(&id)
	return err	
}
