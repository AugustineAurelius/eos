package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/testcontainers/testcontainers-go/modules/clickhouse"

	"github.com/AugustineAurelius/eos/example/common"

	"github.com/AugustineAurelius/eos/example/repository"
	"github.com/AugustineAurelius/eos/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var serviceName = semconv.ServiceNameKey.String("eos-test-repository")

func Test_WithDatabases(t *testing.T) {
	ctx := context.Background()
	conn, err := initConn()
	assert.NoError(t, err)

	res, err := resource.New(ctx, resource.WithAttributes(serviceName))
	assert.NoError(t, err)

	shutdownTracerProvider, err := initTracerProvider(ctx, res, conn)
	assert.NoError(t, err)

	defer func() {
		err = shutdownTracerProvider(ctx)
		assert.NoError(t, err)
	}()
	shutdownMeterProvider, err := initMeterProvider(ctx, res, conn)
	assert.NoError(t, err)

	defer func() {
		err = shutdownMeterProvider(ctx)
		assert.NoError(t, err)
	}()
	name := "go.opentelemetry.io/contrib/examples/otel-collector"
	tracer := otel.Tracer(name)
	meter := otel.Meter(name)
	logger := logger.New(&mode{})

	cases := []struct {
		DatabaseName string
		Provide      func() common.Querier
	}{
		{
			DatabaseName: "sqlite",
			Provide: func() common.Querier {
				db, err := common.NewSqliteInMemory(context.Background(), logger, tracer, meter)
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

				db, err := common.NewPostgres(ctx, common.PgxConnectionProvider{connStr}, logger, tracer, meter)

				assert.NoError(t, err)

				_, err = db.Exec(ctx, `CREATE TABLE if not exists users (id UUID PRIMARY KEY,name TEXT,email TEXT);`)
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
				}, logger, tracer, meter)
				assert.NoError(t, err)

				return &db
			},
		},
		{
			DatabaseName: "clickhouse",
			Provide: func() common.Querier {
				ctx := context.Background()

				user := "clickhouse"
				password := "password"
				dbname := "testdb"

				clickHouseContainer, err := clickhouse.Run(ctx,
					"clickhouse/clickhouse-server:23.3.8.21-alpine",
					clickhouse.WithUsername(user),
					clickhouse.WithPassword(password),
					clickhouse.WithDatabase(dbname),
				)
				assert.NoError(t, err)

				host, err := clickHouseContainer.ConnectionHost(ctx)
				assert.NoError(t, err)

				db, err := common.NewClickhouse(common.ClickhouseConnectionProvider{
					Host:      host,
					User:      user,
					Password:  password,
					Databasse: dbname,
				}, logger, tracer, meter)
				assert.NoError(t, err)
				db.Exec(ctx, `CREATE TABLE users(id UUID, name String, email String) ENGINE = MergeTree() ORDER BY id;`)

				return &db

			},
		},
	}

	for _, c := range cases {
		t.Run(c.DatabaseName, func(t *testing.T) {
			ctx := context.Background()

			db := c.Provide()

			userRepo := repository.New(db)

			id := uuid.New()
			testUser := &repository.User{ID: id, Name: "name", Email: "email"}

			err := userRepo.CreateUser(ctx, testUser)
			assert.NoError(t, err)

			user, err := userRepo.GetUser(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, testUser, user)

			f := repository.NewFilter().AddOneToIDs(id)
			users, err := userRepo.GetManyUsers(ctx, *f)
			assert.NoError(t, err)
			assert.Equal(t, []repository.User{*testUser}, users)

		})
	}
}

func initConn() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient("localhost:4317", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func initMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}

type mode struct {
}

func (m *mode) IsProduction() bool {
	return false
}
