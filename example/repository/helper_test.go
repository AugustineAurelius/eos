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
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
		{ID: uuid.New(), Email: pointerGet(gofakeit.Email()), Balance: gofakeit.Float64()},
	}

	iter := users.All().
		Map(func(u User) User {
			return u
		}).
		Distinct(func(u User) any {
			if u.Email != nil {
				return *u.Email
			}
			return u.Email
		})

	var count int
	for elem := range iter {
		count++
		if count == 2 {
			break
		}
		fmt.Println(elem)
	}

}
