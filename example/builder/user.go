package example_builder

import "github.com/google/uuid"

//go:generate go run github.com/AugustineAurelius/eos generator builder --struct=User --source=user.go --destination=user_builder_gen.go
type User struct {
	name                string
	surname             string
	id                  uuid.UUID
	addesses            []string
	addressesID         []uuid.UUID
	inner               InnerUser
	pointer             *string
	pointerID           *uuid.UUID
	pointerSliceID      *[]uuid.UUID
	pointerSlicePointer *[]string
}

type InnerUser struct {
	Name string
}
