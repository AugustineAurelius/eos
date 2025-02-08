//Code generated by generator, DO NOT EDIT.
package repository


import (
	"context"
	"fmt"
	txrunner "github.com/AugustineAurelius/eos/example/tx_runner" 
  common "github.com/AugustineAurelius/eos/example/common"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)




// UpdateUser updates an existing User in the database.
func (r *Repository) Update(ctx context.Context,  id uuid.UUID, u UserUpdate) error {
	if tx, ok := txrunner.FromContex(ctx); ok {
		return update(ctx, tx, id, u)
    } else {
		return update(ctx, r.db, id, u)
    }
}

func update(ctx context.Context, run common.Querier, id uuid.UUID, u UserUpdate) error {
    b:= sq.Update(TableUser).PlaceholderFormat(sq.Question).Where(sq.Eq{ColumnUserID: id})
    b = ApplySet(b, u)
	query, args := b.MustSql()
	if _, err := run.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to exec update query %s with args %v", query, args)
	}
	return nil 
}

// UpdateUser updates an existing User in the database.
func (r *Repository) UpdateMany(ctx context.Context,  f UserFilter, u  UserUpdate) error {
	if tx, ok := txrunner.FromContex(ctx); ok {
		return updateMany(ctx, tx, f, u)
    } else {
		return updateMany(ctx, r.db, f, u)
    }
}

func updateMany(ctx context.Context, run common.Querier,  f UserFilter, u  UserUpdate) error {
    b:= sq.Update(TableUser).PlaceholderFormat(sq.Question)

	b = ApplyWhere(b, f)

    b = ApplySet(b, u)
	
	query, args := b.MustSql()
	if _, err := run.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to exec update query %s with args %v error = %w", query, args, err)
	}
	return nil 
}



