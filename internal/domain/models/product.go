package models

import (
	"time"

	"github.com/google/uuid"
)

type Type string

func (t Type) IsValid() bool {
	return t == "электроника" || t == "одежда" || t == "обувь"
}

type Product struct {
	ID uuid.UUID
	DateTime *time.Time
	ReceptionID uuid.UUID
	Type Type
}