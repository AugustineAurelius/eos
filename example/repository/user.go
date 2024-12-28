package repository

import "github.com/google/uuid"

//go:generate go run github.com/AugustineAurelius/eos/ generator repository  --type User

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}
