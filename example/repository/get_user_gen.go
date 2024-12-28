//Code generated by generator, DO NOT EDIT.
package repository

import (
	"context"
	"fmt"

    sq	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// GetUser retrieves a User by ID.
func (r *repository) GetUser(ctx context.Context,  id uuid.UUID) (*User, error) {
	query, args := sq.Select(
		ColumnUserID,
		ColumnUserName,
		ColumnUserEmail,
	).
	From(TableUser).
	Where(sq.Eq{ColumnUserID: id}).PlaceholderFormat(sq.Dollar).MustSql()

	var user User
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get query %s with args %v error = %w" , query, args, err)
	}

	return &user, err
}


// GetManyUser retrieves a User by ID.
func (r *repository) GetManyUsers(ctx context.Context, f UserFilter) ([]User, error) {
	b := sq.Select(
		ColumnUserID,
		ColumnUserName,
		ColumnUserEmail,
	).From(TableUser).PlaceholderFormat(sq.Dollar)

	b = ApplyWhere(b, f)

    query, args := 	b.MustSql()

    var users []User

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()


	var user User
	for rows.Next() {
		err := rows.Scan(&user)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
	}


	return users, err
}
