package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string

func (s Status) IsValid() bool {
	return s == "in_progress" || s == "close"
}

type Reception struct {
	ID       uuid.UUID
	DateTime *time.Time
	PVZID    uuid.UUID
	Status   Status
}

type ReceptionFullInfo struct {
	Reception
	Products []Product
}
