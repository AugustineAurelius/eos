// Code generated by generator, DO NOT EDIT.
package common

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type CassandraConnectionProvider struct {
	Hosts    []string
	Port     int
	Keyspace string
	User     string
	Password string
}

func (cp *CassandraConnectionProvider) GetConnectionURL() string {
	return fmt.Sprintf("Contact Points: %v, Port: %d, Keyspace: %s", cp.Hosts, cp.Port, cp.Keyspace)
}

type CassandraDatabase struct {
	session *gocql.Session

	logger logger

	telemetry tracer

	queryCount      int64Counter
	execCount       int64Counter
	queryRowCounter int64Counter
}

func NewCassandraDatabase(cp CassandraConnectionProvider,

	logger logger,

	telemetry tracer,

	metrics metricProvider,

) (CassandraDatabase, error) {
	cluster := gocql.NewCluster(cp.Hosts...)
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: 3,
	}
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cp.User,
		Password: cp.Password,
	}
	cluster.Port = cp.Port
	cluster.Keyspace = cp.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return CassandraDatabase{}, err
	}

	queryCount, err := metrics.Int64Counter("queryCount", metric.WithDescription("Cassandra"))
	if err != nil {
		return CassandraDatabase{}, err
	}
	execCount, err := metrics.Int64Counter("execCount", metric.WithDescription("Cassandra"))
	if err != nil {
		return CassandraDatabase{}, err
	}
	queryRowCounter, err := metrics.Int64Counter("queryRowCounter", metric.WithDescription("Cassandra"))
	if err != nil {
		return CassandraDatabase{}, err
	}

	return CassandraDatabase{
		session: session,

		logger: logger,

		telemetry: telemetry,

		queryCount:      queryCount,
		execCount:       execCount,
		queryRowCounter: queryRowCounter,
	}, nil
}

func (db *CassandraDatabase) Close() error {
	db.session.Close()
	return nil
}

func (db *CassandraDatabase) QueryRow(ctx context.Context, query string, args ...any) row {

	ctx, span := db.telemetry.Start(ctx, "QueryRow", trace.WithAttributes(attribute.String("query", query), attribute.String("db_type", "Cassandra")))
	defer span.End()

	start := time.Now()
	db.logger.Info("Executing QueryRow", zap.String("query", query))

	r := db.session.Query(query, ConvertArgs(args...)...).
		WithContext(context.Background())

	duration := time.Since(start).Seconds()
	db.logger.Info("QueryRow succeeded", zap.Float64("duration", duration))

	db.queryRowCounter.Add(ctx, 1)

	return &CassandraRow{
		Row: r,
	}
}

func (db *CassandraDatabase) Query(ctx context.Context, query string, args ...any) (rows, error) {

	ctx, span := db.telemetry.Start(ctx, "Query", trace.WithAttributes(attribute.String("query", query), attribute.String("db_type", "Cassandra")))
	defer span.End()

	start := time.Now()
	db.logger.Info("Executing Query", zap.String("query", query))

	iter := db.session.Query(RemoveOffset(query), ConvertArgs(args...)...).WithContext(ctx).Iter()

	duration := time.Since(start).Seconds()
	db.logger.Info("Query succeeded", zap.Float64("duration", duration))

	db.queryCount.Add(ctx, 1)

	iter.PageState()
	return &CassandraRows{iter, iter.Scanner()}, nil
}

func (db *CassandraDatabase) Exec(ctx context.Context, query string, args ...any) (result, error) {

	ctx, span := db.telemetry.Start(ctx, "Exec", trace.WithAttributes(attribute.String("query", query), attribute.String("db_type", "Cassandra")))
	defer span.End()

	start := time.Now()
	db.logger.Info("Executing Exec", zap.String("query", query))

	err := db.session.Query(RemoveOffset(query), ConvertArgs(args...)...).
		WithContext(context.Background()).Exec()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		db.logger.Error("Exec failed", zap.Error(err))
		return nil, err
	}

	duration := time.Since(start).Seconds()
	db.logger.Info("Exec succeeded", zap.Float64("duration", duration))

	db.execCount.Add(ctx, 1)

	return &CassandraResult{}, err
}

func (db *CassandraDatabase) Begin(ctx context.Context) (Tx, error) {
	return nil, errors.New("Transactions not supported in Cassandra")
}

type CassandraRows struct {
	Iter    *gocql.Iter
	Scanner gocql.Scanner
}

func (r *CassandraRows) Close() error {
	return r.Iter.Close()
}

func (r *CassandraRows) Err() error {
	return r.Scanner.Err()
}

func (r *CassandraRows) Next() bool {
	return r.Scanner.Next()
}

func (r *CassandraRows) Scan(dest ...any) error {
	return r.Scanner.Scan(dest...)
}

type CassandraRow struct {
	Row *gocql.Query
}

func (r *CassandraRow) Scan(dest ...any) error {
	return r.Row.Scan(dest...)
}

type CassandraResult struct{}

func (r *CassandraResult) RowsAffected() (int64, error) {
	return 0, errors.New("RowsAffected not supported in Cassandra")
}

type CassandraTx struct{}

func (tx *CassandraTx) Query(ctx context.Context, query string, args ...any) (rows, error) {
	return nil, errors.New("Transactions not supported in Cassandra")
}

func (tx *CassandraTx) QueryRow(ctx context.Context, query string, args ...any) row {
	return &CassandraRow{}
}

func (tx *CassandraTx) Exec(ctx context.Context, query string, args ...any) (result, error) {
	return nil, errors.New("Transactions not supported in Cassandra")
}

func (tx *CassandraTx) Commit(ctx context.Context) error {
	return errors.New("Transactions not supported in Cassandra")
}

func (tx *CassandraTx) Rollback(ctx context.Context) error {
	return errors.New("Transactions not supported in Cassandra")
}

func ConvertArgs(args ...any) []any {
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case uuid.UUID:
			args[i] = ConvertUUID(v)
		default:
			continue
		}
	}

	return args
}

func ConvertUUID(u uuid.UUID) gocql.UUID {
	var g gocql.UUID
	copy(g[:], u[:])
	return g
}

func RemoveOffset(query string) string {
	pattern := `\boffset\s+\d+\b`
	re := regexp.MustCompile(pattern)
	modifiedQuery := re.ReplaceAllString(strings.ToLower(query), "")
	modifiedQuery = strings.TrimSpace(modifiedQuery)
	return modifiedQuery
}
