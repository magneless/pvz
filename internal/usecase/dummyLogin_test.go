package usecase_test

import (
	"testing"

	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestDummyLoginUsecase_CreateToken(t *testing.T) {
	uc := usecase.NewDummyLoginUsecase()

	t.Run("valid role", func(t *testing.T) {
		token, err := uc.CreateToken(models.DummyLogin{Role: models.Role("employee")})
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("invalid role", func(t *testing.T) {
		token, err := uc.CreateToken(models.DummyLogin{Role: models.Role("invalid")})
		require.ErrorIs(t, err, usecase.ErrWrongRole)
		require.Empty(t, token)
	})
}
