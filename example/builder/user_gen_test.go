package example_builder

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func Test_UserBuilder(t *testing.T) {
	addressID := uuid.New()
	user := NewUserBuilder().
		Name("user").
		Surname("Surname").
		AddOneToAddressesID(addressID).
		Build()

	expectedUser := User{
		name:        "user",
		surname:     "Surname",
		addressesID: []uuid.UUID{addressID},
	}

	if !reflect.DeepEqual(user, expectedUser) {
		t.Fatal(user, expectedUser)
	}
}
