package repository

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// User represents the User message.
type UserModel struct {
	Id    uuid.UUID
	Name  string
	Email string
}

// Table and column name constants for User
const (
	TableUser       = "users"
	ColumnUserId    = "id"
	ColumnUserName  = "name"
	ColumnUserEmail = "email"
)

// Scan implements the Scanner interface for UserModel.
func (m *UserModel) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	// Assuming the value is a map[string]interface{} representing the row
	row, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected map[string]interface{}, got %T", value)
	}

	// Scan each field
	if val, ok := row[ColumnUserId]; ok {
		switch v := val.(type) {
		case uuid.UUID:
			m.Id = v
		default:
			return fmt.Errorf("unexpected type for Id: got %T, expected uuid.UUID", val)
		}
	}
	if val, ok := row[ColumnUserName]; ok {
		switch v := val.(type) {
		case string:
			m.Name = v
		default:
			return fmt.Errorf("unexpected type for Name: got %T, expected string", val)
		}
	}
	if val, ok := row[ColumnUserEmail]; ok {
		switch v := val.(type) {
		case string:
			m.Email = v
		default:
			return fmt.Errorf("unexpected type for Email: got %T, expected string", val)
		}
	}

	return nil
}

// Value implements the Valuer interface for UserModel.
func (m UserModel) Value() (driver.Value, error) {
	// Convert the struct into a map[string]interface{}
	row := make(map[string]interface{})
	row[ColumnUserId] = m.Id
	row[ColumnUserName] = m.Name
	row[ColumnUserEmail] = m.Email
	return row, nil
}
