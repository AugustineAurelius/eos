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
	_, err := tx.ExecContext(ctx, `CREATE TABLE tasks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR NOT NULL,
		description VARCHAR,
		created_by VARCHAR NOT NULL,
		doer VARCHAR NOT NULL,
		done BOOLEAN NOT NULL DEFAULT false,
		repeatable BOOLEAN NOT NULL DEFAULT false,
		repeat_after INTEGER,
		do_before TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ
	);

	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER update_task_updated_at
	BEFORE UPDATE ON tasks
	FOR EACH ROW
	EXECUTE FUNCTION update_updated_at_column();`)
	return err
}

func downInitMigrationForfuufu(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TRIGGER IF EXISTS update_task_updated_at ON tasks;
	DROP TABLE IF EXISTS tasks;`)
	return err
}
