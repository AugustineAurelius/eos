package repository

//go:generate go run github.com/AugustineAurelius/eos/ generator repository  --type User  --default_id=true --table=users
type User struct {
	ID    int
	Name  string
	Email *string
	// Booler  bool
	Balance float64
	// Balance2 decimal.Decimal
	// Created time.Time
	// Addresses []string
	// UserTime int
}
