package repository

//go:generate go run github.com/AugustineAurelius/eos/ generator repository  --type User  --default_id=true
type User struct {
	ID    int
	Name  string
	Email *string
	// Booler  bool
	Balance float64
	// Created time.Time
	// Addresses []string
	// UserTime int
}
