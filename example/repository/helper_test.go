package repository

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func pointerGet[T any](t T) *T {
	return &t
}
func Test_Iter(t *testing.T) {

	users := Users{
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email())},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email())},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email())},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email())},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email())},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email())},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email())},
		{ID: uuid.New(), Email: pointerGet("gofakeit.Email()")},
		{ID: uuid.New(), Email: pointerGet("gofakeit.Email()")},
	}

	users.All().
		Map(func(u User) User {
			return u
		}).
		// FilterByEmail(&email).
		Distinct(func(u User) any {
			if u.Email != nil {
				return *u.Email
			}
			return u.Email
		}).
		ForEach(func(u User) {
			fmt.Println(*u.Email)
		})

}
