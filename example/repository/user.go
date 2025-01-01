package repository

import (
	"github.com/google/uuid"
)

//go:generate go run github.com/AugustineAurelius/eos/ generator repository  --type User --with_tx=true --tx_path=example/tx_runner

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}
