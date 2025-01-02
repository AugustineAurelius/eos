//Code generated by generator, DO NOT EDIT.
package txrunner

import (
	"context"
	"fmt"
	common "github.com/AugustineAurelius/eos/example/common"
)

var TxKey = struct{}{}

type txRunner struct {
	db common.Begginer
}

func New(db common.Begginer) *txRunner {
	return &txRunner{db}
}

func (tr *txRunner) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := tr.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic handled")
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				// add logger here
				fmt.Println(rollbackErr)
			}	
		} else if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				// add logger here
				fmt.Println(rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				fmt.Println(commitErr)
			}
		}
	}()
	ctxWithTx := context.WithValue(ctx, TxKey, tx)
	err = fn(ctxWithTx)
	return err
}

func FromContex(ctx context.Context) (common.Tx, bool) {
		tx, ok := ctx.Value(TxKey).(common.Tx)
	if !ok {
		return nil, false
	}
	return tx, true
}
