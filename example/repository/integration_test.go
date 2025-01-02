package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/AugustineAurelius/eos/example/common"
	"github.com/AugustineAurelius/eos/example/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_WithDatabases(t *testing.T) {

	cases := []struct {
		DatabaseName string
		Provide      func() common.Querier
	}{
		{
			DatabaseName: "sqlite",
			Provide: func() common.Querier {
				db, err := common.NewSqliteInMemory(context.Background())
				assert.NoError(t, err)

				_, err = db.Exec(context.Background(), `CREATE TABLE users (id TEXT PRIMARY KEY, name TEXT, email TEXT);`)
				assert.NoError(t, err)

				return &db

			},
		},
		{
			DatabaseName: "postgres",
			Provide: func() common.Querier {
				ctx := context.Background()

				c, err := postgres.RunContainer(ctx,
					testcontainers.WithImage("postgres:15.3-alpine"),
					postgres.WithDatabase("users-db"),
					postgres.WithUsername("postgres"),
					postgres.WithPassword("postgres"),
					testcontainers.WithWaitStrategy(
						wait.ForLog("database system is ready to accept connections").
							WithOccurrence(2).WithStartupTimeout(5*time.Second)),
				)
				assert.NoError(t, err)

				connStr, err := c.ConnectionString(ctx, "sslmode=disable")
				assert.NoError(t, err)

				db, err := common.NewPostgres(ctx, common.PgxConnectionProvider{connStr})
				assert.NoError(t, err)

				_, err = db.Exec(ctx, `CREATE TABLE users (  id UUID PRIMARY KEY,name TEXT,email TEXT);`)
				assert.NoError(t, err)

				return &db
			},
		},
		{
			DatabaseName: "cassandra",
			Provide: func() common.Querier {
				ctx := context.Background()

				c, err := cassandra.Run(ctx, "cassandra:4.1.3", cassandra.WithInitScripts("./cassandra.cql"),
					testcontainers.WithEnv(map[string]string{
						"CASSANDRA_HOST":     "cassandra",
						"CASSANDRA_USER":     "user",
						"CASSANDRA_PASSWORD": "pass",
					}))
				assert.NoError(t, err)

				host, err := c.ConnectionHost(ctx)
				assert.NoError(t, err)

				db, err := common.NewCassandraDatabase(common.CassandraConnectionProvider{
					Hosts:    []string{host},
					Port:     9042,
					User:     "user",
					Password: "pass",
					Keyspace: "test",
				})
				assert.NoError(t, err)

				return &db
			},
		},
	}

	for _, c := range cases {
		t.Run(c.DatabaseName, func(t *testing.T) {
			db := c.Provide()

			userRepo := repository.New(db)

			id := uuid.New()
			testUser := &repository.User{ID: id, Name: "name", Email: "email"}

			err := userRepo.CreateUser(context.Background(), testUser)
			assert.NoError(t, err)

			user, err := userRepo.GetUser(context.Background(), id)
			assert.NoError(t, err)
			assert.Equal(t, testUser, user)

			f := repository.NewFilter().AddOneToIDs(id)
			users, err := userRepo.GetManyUsers(context.Background(), *f)
			assert.NoError(t, err)
			assert.Equal(t, []repository.User{*testUser}, users)

		})
	}
}
