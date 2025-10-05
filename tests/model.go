package main

import (
	"time"

	"github.com/shopspring/decimal"
)

//go:generate go tool eos generator repository  --type NavHistory  --default_id=true --table=nav_history
type NavHistory struct {
	ID        int             `db:"id"`
	RobotID   int64           `db:"robot_id"`
	Nav       decimal.Decimal `db:"nav"`
	CreatedAt time.Time       `db:"created_at"`
}
