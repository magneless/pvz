package usecase

import (
	"errors"

	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/pkg/jwt"
)

var ErrWrongRole = errors.New("wrong role")

type DummyLoginUsecase struct{}

func NewDummyLoginUsecase() *DummyLoginUsecase {
	return &DummyLoginUsecase{}
}

func (*DummyLoginUsecase) CreateToken(dummyLogin models.DummyLogin) (string, error) {
	if !dummyLogin.Role.IsValid() {
		return "", ErrWrongRole
	}
	return jwt.GenerateAccessToken(string(dummyLogin.Role))
}
