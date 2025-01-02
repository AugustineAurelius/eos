//Code generated by generator, DO NOT EDIT.
package common

import (
	"context"
	"regexp"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConnectionProvider struct {
	URL string
}

func (p *PgxConnectionProvider) GetConnectionURL() string {
	return p.URL
}

type PgxPoolDB struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, provider PgxConnectionProvider) (PgxPoolDB, error){
	url := provider.GetConnectionURL()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return PgxPoolDB{}, err
	}
	return PgxPoolDB{pool}, nil
}

func (p *PgxPoolDB) Close() error {
	p.pool.Close()
	return nil
}

func (p *PgxPoolDB) Query(ctx context.Context, query string, args ...any) (rows, error) {
	rows, err := p.pool.Query(ctx, ReplaceQuestions(query), args...)
	if err != nil {
		return nil, err
	}
	return &PgxRows{rows}, nil
}

func (p *PgxPoolDB) QueryRow(ctx context.Context, query string, args ...any) row {
	return p.pool.QueryRow(ctx, ReplaceQuestions(query), args...)
}

func (p *PgxPoolDB) Exec(ctx context.Context, query string, args ...any) (result, error) {
	r, err := p.pool.Exec(ctx, ReplaceQuestions(query), args...)
	return &PgxResult{r}, err
}

func (p *PgxPoolDB) BeginTransaction(ctx context.Context) (Tx, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PgxTx{tx}, nil
}

type PgxRows struct {
	pgx.Rows
}

func (p *PgxRows) Next() bool {
	return p.Rows.Next()
}

func (p *PgxRows) Scan(dest ...any) error {
	return p.Rows.Scan(dest...)
}

func (p *PgxRows) Close() error {
	p.Rows.Close()
	return nil
}

func (p *PgxRows) Err() error {
	return p.Rows.Err()
}

type PgxTx struct {
	pgx.Tx
}

func (p *PgxTx) Query(ctx context.Context, query string, args ...any) (rows, error) {
	rows, err := p.Tx.Query(ctx, ReplaceQuestions(query), args...)
	if err != nil {
		return nil, err
	}
	return &PgxRows{rows}, nil
}

func (p *PgxTx) QueryRow(ctx context.Context, query string, args ...any) row {
	return p.Tx.QueryRow(ctx, ReplaceQuestions(query), args...)
}

func (p *PgxTx) Exec(ctx context.Context, query string, args ...any) (result, error) {
	r, err := p.Tx.Exec(ctx, ReplaceQuestions(query), args...)
	return &PgxResult{r}, err
}

func (p *PgxTx) Commit(ctx context.Context) error {
	return p.Tx.Commit(ctx)
}

func (p *PgxTx) Rollback(ctx context.Context) error {
	return p.Tx.Rollback(ctx)
}

type PgxResult struct {
	pgconn.CommandTag
}

func (r *PgxResult) RowsAffected() (int64, error) {
	return r.CommandTag.RowsAffected(), nil
}


func ReplaceQuestions(query string) string {
    var count int
    re := regexp.MustCompile(`\?`)
    return re.ReplaceAllStringFunc(query, func(s string) string {
        count++
        return "$" + strconv.Itoa(count)
    })
}
