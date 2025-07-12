package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Preferences Preferences
	Orders      []Order
	CreatedAt   time.Time
	ID          uuid.UUID
	Name        string
	Email       string
	Address     Address
}
