package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UpdateUser updates an existing User in the database.
func UpdateUser(ctx context.Context, db *pgxpool.Pool, user *User) error {
	query, args, err := sq.Update(TableUser).
		SetMap(map[string]interface{}{
			ColumnUserId: user.Id,
			ColumnUserName: user.Name,
			ColumnUserEmail: user.Email,
		}).
		Where(sq.Eq{ColumnUserId: user.Id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = db.Exec(ctx, query, args...)
	return err
}
