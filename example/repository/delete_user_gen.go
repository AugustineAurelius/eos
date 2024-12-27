package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

// DeleteUser deletes a User by ID.
func DeleteUser(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) error {
	query, args, err := sq.Delete(TableUser).
		Where(sq.Eq{ColumnUserId: id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = db.Exec(ctx, query, args...)
	return err
}
