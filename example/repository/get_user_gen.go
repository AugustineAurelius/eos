//Code generated by generator, DO NOT EDIT.
package repository


import (
	"context"
	"fmt"

	txrunner "github.com/AugustineAurelius/eos/example/tx_runner" 



	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)



// GetUser retrieves a User by ID.
func (r *repository) GetUser(ctx context.Context,  id uuid.UUID) (*User, error) {
	if tx, ok := txrunner.FromContex(ctx); ok {
		return get(ctx, tx, id)
    } else {
		return get(ctx, r.db, id)
    }
}


func get(ctx context.Context, run runner, id uuid.UUID) (*User, error){
	query, args := sq.Select(
		ColumnUserID,
		ColumnUserName,
		ColumnUserEmail,
	).
	From(TableUser).
	Where(sq.Eq{ColumnUserID: id}).PlaceholderFormat(sq.Dollar).MustSql()

	var user User
	err := run.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get query %s with args %v error = %w" , query, args, err)
	}

	return &user, err

}

// GetManyUser retrieves a User by filter.
func (r *repository) GetManyUsers(ctx context.Context, f UserFilter) ([]User, error) {
	if tx, ok := txrunner.FromContex(ctx); ok {
		return getMany(ctx, tx, f)
    } else {
		return getMany(ctx, r.db, f)
    }
}

func getMany(ctx context.Context, run runner, f UserFilter) ([]User, error) {
	b := sq.Select(
		ColumnUserID,
		ColumnUserName,
		ColumnUserEmail,
	).From(TableUser).PlaceholderFormat(sq.Dollar)

	b = ApplyWhere(b, f)

    query, args := 	b.MustSql()

    var users []User

	rows, err := run.Query(ctx, query, args...)
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



