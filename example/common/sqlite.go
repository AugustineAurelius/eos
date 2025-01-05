//Code generated by generator, DO NOT EDIT.
package common

import (
	"context"
	"time"
    "database/sql"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	
	"go.opentelemetry.io/otel/metric"
	
	
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
    
)

type SQLiteConnectionProvider struct {
    URL string
}

func (s *SQLiteConnectionProvider) GetConnectionURL() string {
    return s.URL
}

type SQLiteDB struct {
    db *sql.DB
	
	logger logger
	
	
	telemetry tracer
	
    
	queryCount      int64Counter
	execCount       int64Counter
	queryRowCounter int64Counter
	
}

func NewSQLite(ctx context.Context, provider SQLiteConnectionProvider,
	
	logger logger,
	
	
	telemetry tracer, 
	
	
	metrics metricProvider,
	
) (SQLiteDB, error){
    url := provider.GetConnectionURL()
    db, err := sql.Open("sqlite3", url)
    if err != nil {
        return SQLiteDB{}, err
    }
	
	queryCount, err := metrics.Int64Counter("queryCount",metric.WithDescription("SQLite"))
	if err != nil {
		return SQLiteDB{}, err
	}
	execCount, err := metrics.Int64Counter("execCount",metric.WithDescription("SQLite"))
	if err != nil {
		return SQLiteDB{}, err
	}
	queryRowCounter, err := metrics.Int64Counter("queryRowCounter",metric.WithDescription("SQLite"))
	if err != nil {
		return SQLiteDB{}, err 
	}
	
    return SQLiteDB{db: db,
		
		logger:          logger,
		
		
		telemetry:       telemetry,
		
		
		queryCount:      queryCount,
		execCount:       execCount,
		queryRowCounter: queryRowCounter,
		
	}, nil
}

func (s *SQLiteDB) Close() error {
    return s.db.Close()
}

func NewSQLiteInMemory(ctx context.Context,
	
	logger logger,
	
	
	telemetry tracer, 
	
	
	metrics metricProvider,
	
) (SQLiteDB, error){
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        return SQLiteDB{}, err
    }
	
	queryCount, err := metrics.Int64Counter("queryCount",metric.WithDescription("SQLite"))
	if err != nil {
		return SQLiteDB{}, err
	}
	execCount, err := metrics.Int64Counter("execCount",metric.WithDescription("SQLite"))
	if err != nil {
		return SQLiteDB{}, err
	}
	queryRowCounter, err := metrics.Int64Counter("queryRowCounter",metric.WithDescription("SQLite"))
	if err != nil {
		return SQLiteDB{}, err 
	}
	
    return SQLiteDB{db: db,
		
		logger:          logger,
		
		
		telemetry:       telemetry,
		
		
		queryCount:      queryCount,
		execCount:       execCount,
		queryRowCounter: queryRowCounter,
		
}, nil
}

func (db *SQLiteDB) QueryRow(ctx context.Context, query string, args ...any) row {
	
	ctx, span := db.telemetry.Start(ctx, "QueryRow",trace.WithAttributes(attribute.String("query", query), attribute.String("db_type", "SQLite")))
	defer span.End()
    
	
	start := time.Now()
	db.logger.Info("Executing QueryRow", zap.String("query", query))
	
	row := db.db.QueryRowContext(ctx, query, args...)
	
	duration := time.Since(start).Seconds()
	db.logger.Info("QueryRow succeeded", zap.Float64("duration", duration))
	
	
	db.queryRowCounter.Add(ctx, 1)
	
    return row
}

func (db *SQLiteDB) Query(ctx context.Context, query string, args ...any) (rows, error) {
	
	ctx, span := db.telemetry.Start(ctx, "Query",trace.WithAttributes(attribute.String("query", query), attribute.String("db_type", "SQLite")))
	defer span.End()
    
	
    start := time.Now()
    db.logger.Info("Executing Query", zap.String("query", query))
	
    rows, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
        
		
        db.logger.Error("Query failed", zap.Error(err))
		
		return nil, err
	}
	
	duration := time.Since(start).Seconds()
    db.logger.Info("Query succeeded", zap.Float64("duration", duration))
	
	
	db.queryCount.Add(ctx, 1)
    
    return &SQLiteRows{rows}, nil
}

func (db *SQLiteDB) Exec(ctx context.Context, query string, args ...any) (result, error) {
	
    ctx, span := db.telemetry.Start(ctx, "Exec",trace.WithAttributes(attribute.String("query", query), attribute.String("db_type", "SQLite")))
    defer span.End()
    
	
	start := time.Now()
    db.logger.Info("Executing Exec", zap.String("query", query))
	
    r, err := db.db.ExecContext(ctx, query, args...)
    if err != nil {
        
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
        
		
        db.logger.Error("Exec failed", zap.Error(err))
		
		return nil, err
    } 
	
    duration := time.Since(start).Seconds()
    db.logger.Info("Exec succeeded", zap.Float64("duration", duration))
	
	
	db.execCount.Add(ctx, 1)
    
    return r, err
}

func (s *SQLiteDB) BeginTransaction(ctx context.Context) (Tx, error) {
    tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
    if err != nil {
        return nil, err
    }
    return &SQLiteTx{tx}, nil
}

type SQLiteRows struct {
    *sql.Rows
}

func (s *SQLiteRows) Next() bool {
    return s.Rows.Next()
}

func (s *SQLiteRows) Scan(dest ...any) error {
    return s.Rows.Scan(dest...)
}

func (s *SQLiteRows) Close() error {
    return s.Rows.Close()
}

func (s *SQLiteRows) Err() error {
    return s.Rows.Err()
}

type SQLiteTx struct {
	*sql.Tx
}

func (s *SQLiteTx) Query(ctx context.Context, query string, args ...any) (rows, error) {
	rows, err := s.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLiteRows{rows}, nil
}

func (s *SQLiteTx) QueryRow(ctx context.Context, query string, args ...any) row {
	return s.Tx.QueryRowContext(ctx, query, args...)
}

func (s *SQLiteTx) Exec(ctx context.Context, query string, args ...any) (result, error) {
	return s.Tx.ExecContext(ctx, query, args...)
}

func (s *SQLiteTx) Commit(ctx context.Context) error {
	return s.Tx.Commit()
}

func (s *SQLiteTx) Rollback(ctx context.Context) error {
	return s.Tx.Rollback()
}
