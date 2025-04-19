package repository

import (
	"github.com/google/uuid"
)

//go:generate go run github.com/AugustineAurelius/eos/ generator repository  --type User  --tx_path=example/tx_runner --common_path=example/common

type User struct {
	ID    uuid.UUID
	Name  string
	Email *string
	// Booler  bool
	// Balance float64
	// Created time.Time
	// Addresses []string
	// UserTime int
}
