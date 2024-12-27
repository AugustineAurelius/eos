package repository

import (
	"context"
	"fmt"

    sq	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

// GetUser retrieves a User by ID.
func GetUser(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (*User, error) {
	query, args, err := sq.Select(
		ColumnUserId,
		ColumnUserName,
		ColumnUserEmail,
	).
	From(TableUser).
	Where(sq.Eq{ColumnUserId: id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	var user User
	err = db.QueryRow(ctx, query, args...).Scan(&{{$.MessageName | lower}}.Id, &{{$.MessageName | lower}}.Name, &{{$.MessageName | lower}}.Email)
	return &user, err
}
