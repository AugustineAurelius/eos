package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateUser inserts a new User into the database.
func CreateUser(ctx context.Context, db *pgxpool.Pool, user *User) error {
	query, args, err := sq.Insert(TableUser).
		Columns(ColumnUserId, ColumnUserName, ColumnUserEmail).
		Values(user.Id, user.Name, user.Email).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = db.Exec(ctx, query, args...)
	return err
}
