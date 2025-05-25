package migration

import (
	"context"
	"errors"
	"fmt"
	"math"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type DNSer interface {
	DNS() string
}

const (
	minVersion = int64(0)
	maxVersion = math.MaxInt64
)

func CheckMigrations(ctx context.Context, postgresConfig DNSer) error {
	db, err := goose.OpenDBWithDriver("postgres", postgresConfig.DNS())
	if err != nil {
		return err
	}
	defer db.Close()

	version, err := goose.GetDBVersionContext(ctx, db)
	if err != nil {
		return err
	}

	if version == 0{
		return errors.New("should run fuufu migrate up")
	}

	var current int64
	migrations, err := goose.CollectMigrations(".", minVersion, maxVersion)
	if err != nil {
		return fmt.Errorf("failed to collect migrations: %w", err)
	}
	if len(migrations) > 0 {
		current = migrations[len(migrations)-1].Version
	}

	if version != current {
		return errors.New("not all migrations applied should run fuufu migrate")
	}

	return nil
}
