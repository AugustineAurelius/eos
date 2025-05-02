package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitMigrationForfuufu, downInitMigrationForfuufu)
}

func upInitMigrationForfuufu(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, ``)
	return err
}

func downInitMigrationForfuufu(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, ``)
	return err
}
