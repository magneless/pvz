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

func TestPVZRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := postgresql.NewPVZRepository(db)
	ctx := context.Background()

	id := uuid.New()
	date := time.Now()
	pvz := models.PVZ{
		ID:               id,
		RegistrationDate: &date,
		City:             "Moscow",
	}

	mock.ExpectExec("INSERT INTO pvz").
		WithArgs(id, date, pvz.City).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(ctx, pvz)
	require.NoError(t, err)
}
