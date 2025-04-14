package postgresql_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/repository/postgresql"
	"github.com/stretchr/testify/require"
)

func ptTime(t time.Time) *time.Time {
	return &t
}

func TestReceptionRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := postgresql.NewReceptionRepository(db)
	ctx := context.Background()

	reception := models.Reception{
		ID:       uuid.New(),
		DateTime: ptTime(time.Now()),
		Status:   "in_progress",
		PVZID:    uuid.New(),
	}

	mock.ExpectQuery("INSERT INTO receptions").
		WithArgs(reception.ID, *reception.DateTime, reception.Status, reception.PVZID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(reception.ID))

	err := repo.Create(ctx, reception)
	require.NoError(t, err)
}

func TestReceptionRepository_CloseReception(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := postgresql.NewReceptionRepository(db)
	ctx := context.Background()
	pvzID := uuid.New()
	closedID := uuid.New()

	mock.ExpectQuery("UPDATE receptions").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(closedID))

	err := repo.CloseReception(ctx, pvzID)
	require.NoError(t, err)
}
