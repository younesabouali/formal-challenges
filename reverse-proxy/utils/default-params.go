package utils

import (
	"time"

	"github.com/google/uuid"
)

type DefaultParams struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetDefaultParams() DefaultParams {
	return DefaultParams{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
