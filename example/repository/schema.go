package repository

import "github.com/google/uuid"

// User represents the User message.
type UserModel struct {
	Id uuid.UUID
	Name string
	Email string
}

// Table and column name constants for User
const (
	TableUser = "users"
	ColumnUserId = "id"
	ColumnUserName = "name"
	ColumnUserEmail = "email"
)
