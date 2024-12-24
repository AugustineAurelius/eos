package example_builder

import "github.com/google/uuid"

//go:generate eos generator builder --struct=User --source=user.go --destination=user_builder.go
type User struct {
	name        string
	surname     string
	id          uuid.UUID
	addesses    []string
	addressesID []uuid.UUID
	inner       InnerUser
}

type InnerUser struct {
	Name string
}
