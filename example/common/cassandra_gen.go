//Code generated by generator, DO NOT EDIT.
package common

import (
	"context"
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
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
}

func NewCassandraDatabase(cp CassandraConnectionProvider) (CassandraDatabase, error) {
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
	return CassandraDatabase{
		session: session,
	}, nil
}

func (db *CassandraDatabase) Close() error {
	db.session.Close()
	return nil
}

func (db *CassandraDatabase) Query(ctx context.Context, query string, args ...any) (rows, error) {
	iter := db.session.Query(query, ConvertArgs(args...)...).WithContext(ctx).Iter()
	return &CassandraRows{iter, iter.Scanner()}, nil
}

func (db *CassandraDatabase) QueryRow(ctx context.Context, query string, args ...any) row {
	r := db.session.Query(query, ConvertArgs(args...)...).
		WithContext(context.Background())
	return &CassandraRow{
		Row: r,
	}
}

func (db *CassandraDatabase) Exec(ctx context.Context, query string, args ...any) (result, error) {
	err := db.session.Query(query, ConvertArgs(args...)...).
		WithContext(context.Background()).Exec()
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
