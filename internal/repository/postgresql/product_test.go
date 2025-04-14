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

func TestProductRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := postgresql.NewProductRepository(db)
	ctx := context.Background()
	product := models.Product{
		ID:       uuid.New(),
		DateTime: ptrTime(time.Now()),
		Type:     "tv",
	}
	pvzID := uuid.New()
	receptionID := uuid.New()

	mock.ExpectQuery("INSERT INTO products").
		WithArgs(product.ID, *product.DateTime, product.Type, pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"reception_id"}).AddRow(receptionID))

	gotID, err := repo.Create(ctx, product, pvzID)
	require.NoError(t, err)
	require.Equal(t, receptionID, gotID)
}

func TestProductRepository_DeleteProduct(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := postgresql.NewProductRepository(db)
	ctx := context.Background()
	pvzID := uuid.New()
	deletedID := uuid.New()

	mock.ExpectQuery("DELETE FROM products").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(deletedID))

	err := repo.DeleteProduct(ctx, pvzID)
	require.NoError(t, err)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
