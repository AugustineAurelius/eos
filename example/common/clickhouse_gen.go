//Code generated by generator, DO NOT EDIT.
package common

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickhouseConnectionProvider struct {
	Host      string
	Port      int
	Databasse string
	User      string
	Password  string
}

type ClickHouseQuerier struct {
	client *sql.DB
}

func NewClickhouse(cp ClickhouseConnectionProvider) (ClickHouseQuerier, error) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{cp.Host},
		Auth: clickhouse.Auth{
			Database: cp.Databasse,
			Username: cp.User,
			Password: cp.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                true,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	return ClickHouseQuerier{conn}, nil
}

func (c *ClickHouseQuerier) Query(ctx context.Context, query string, args ...any) (rows, error) {
	rows, err := c.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &ClickHouseRows{rows}, nil
}

func (c *ClickHouseQuerier) QueryRow(ctx context.Context, query string, args ...any) row {
	row := c.client.QueryRowContext(ctx, query, args...)
	return &ClickHouseRow{row}
}

func (c *ClickHouseQuerier) Exec(ctx context.Context, query string, args ...any) (result, error) {
	res, err := c.client.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &ClickHouseResult{res}, nil
}

type ClickHouseRows struct {
	rows
}

func (r *ClickHouseRows) Err() error {
	return r.rows.Err()
}

func (r *ClickHouseRows) Next() bool {
	return r.rows.Next()
}

func (r *ClickHouseRows) Close() error {
	return r.rows.Close()
}

func (r *ClickHouseRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

type ClickHouseRow struct {
	row
}

func (r *ClickHouseRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

type ClickHouseResult struct {
	result
}

func (r *ClickHouseResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}

type ClickHouseBegginer struct {
	*ClickHouseQuerier
}

func (b *ClickHouseBegginer) Begin(ctx context.Context) (Tx, error) {
	return &ClickHouseTx{}, errors.New("not supported")
}

type ClickHouseTx struct {
	*ClickHouseQuerier
}

func (t *ClickHouseTx) Commit(ctx context.Context) error {
	return errors.New("not supported")

}

func (t *ClickHouseTx) Rollback(ctx context.Context) error {
	return errors.New("not supported")
}
