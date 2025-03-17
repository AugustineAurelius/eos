//Code generated by generator, DO NOT EDIT.
package repository


import (
	"context"
	"fmt"
	txrunner "github.com/AugustineAurelius/eos/example/tx_runner" 
  common "github.com/AugustineAurelius/eos/example/common"

	sq "github.com/Masterminds/squirrel"
)




// UpdateUser updates an existing User in the database.
func (r *Repository) Update(ctx context.Context, u Update, opts ...FilterOpt) error {
	if tx, ok := txrunner.FromContex(ctx); ok {
		return update(ctx, tx, u, opts...)
    } else {
		return update(ctx, r.db, u, opts...)
    }
}

func update(ctx context.Context, run common.Querier, u Update, opts ...FilterOpt) error {
    b:= sq.Update(TableUser).PlaceholderFormat(sq.Question)
	f := &Filter{}
	for i := 0; i < len(opts); i++ {
		opts[i](f)
	}
	b = ApplyWhere(b, *f)
    b = ApplySet(b, u)
	query, args := b.MustSql()
	if _, err := run.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to exec update query %s with args %v", query, args)
	}
	return nil 
}




