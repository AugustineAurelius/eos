//Code generated by generator, DO NOT EDIT.
package repository



import (
	"context"
	"fmt"
	txrunner "github.com/AugustineAurelius/eos/example/tx_runner" 
	common "github.com/AugustineAurelius/eos/example/common"

	sq "github.com/Masterminds/squirrel"
)



// CreateUser inserts a new User into the database.
func (r *repository) Create (ctx context.Context, user *User) error {
	if tx, ok := txrunner.FromContex(ctx); ok {
		return create(ctx, tx, user)
    } else {
		return create(ctx, r.db, user)
    }
}


func create(ctx context.Context, run common.Querier, user *User) error {
	model := Converter(*user)
	query, args := sq.Insert(TableUser).
		Columns(ColumnUserID, ColumnUserName, ColumnUserEmail).
		Values(model.Values()...).PlaceholderFormat(sq.Question).MustSql()

	if _, err := run.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to exec create query %s with args %v error = %w", query, args, err)
	}
	return nil
}

