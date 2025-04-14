package models

import (
	"time"

	"github.com/google/uuid"
)

type City string

func (c City) IsValid() bool {
	return c == "Москва" || c == "Санкт-Петербург" || c == "Казань"
}

type PVZ struct {
	ID               uuid.UUID
	RegistrationDate *time.Time
	City             City
}

type PVZFullInfo struct {
	PVZ
	Receptions []ReceptionFullInfo
}
