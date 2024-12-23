package example

import "github.com/google/uuid"

//go:generate go run github.com/AugustineAurelius/eos generator builder --struct=User --source=user.go --destination=user_builder.go
type User struct {
	Name        string
	Surname     string
	ID          uuid.UUID
	Addesses    []string
	AddressesID []uuid.UUID
}
