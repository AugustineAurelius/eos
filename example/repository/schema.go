// Code generated by generator, DO NOT EDIT.
package repository

import (
	"github.com/google/uuid"
)

// Table and column name constants for User
const (
	TableUser          = "users"
	ColumnUserID       = "id"
	ColumnUserName     = "name"
	ColumnUserEmail    = "email"
	ColumnUserUserTime = "user_time"
)

// User represents the User message.
type UserModel struct {
	ID string

	Name string

	Email string

	UserTime int
}

func (m UserModel) Values() []any {
	return []any{
		m.ID,
		m.Name,
		m.Email,
		m.UserTime,
	}
}

func Converter(user User) UserModel {
	return UserModel{

		ID: user.ID.String(),

		Name: user.Name,

		Email: user.Email,

		UserTime: user.UserTime,
	}
}

func ReverseConverter(userModel UserModel) User {
	return User{

		ID: uuid.MustParse(userModel.ID),

		Name: userModel.Name,

		Email: userModel.Email,

		UserTime: userModel.UserTime,
	}
}
