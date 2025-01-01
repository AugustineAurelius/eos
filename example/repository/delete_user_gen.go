//Code generated by generator, DO NOT EDIT.
package repository


import (
	"context"
	"fmt"
	txrunner "github.com/AugustineAurelius/eos/example/tx_runner" 

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)



// DeleteUser deletes a User by ID.
func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if tx, ok := txrunner.FromContex(ctx); ok {
		return delete(ctx, tx, id)
    } else {
		return delete(ctx, r.db, id)
    }
}

func delete(ctx context.Context, run runner,id uuid.UUID) error {
	query, args := sq.Delete(TableUser).
		Where(sq.Eq{ColumnUserID: id}).PlaceholderFormat(sq.Dollar).MustSql()

	if _, err := run.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to exec delete query %s with args %v error = %w", query, args, err)
	}
	return nil
}

// DeleteManyUser retrieves a User by filter.
func (r *repository) DeleteManyUser(ctx context.Context, f UserFilter) error {
	if tx, ok := txrunner.FromContex(ctx); ok {
		return deleteMany(ctx, tx, f)
    } else {
		return deleteMany(ctx, r.db, f)
    }
}

func deleteMany(ctx context.Context,run runner, f UserFilter) error {
	b := sq.Delete(TableUser).PlaceholderFormat(sq.Dollar)

	b = ApplyWhere(b, f)

    query, args := 	b.MustSql()

	_, err := run.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error querying database: %w", err)
	}

	return err
}



