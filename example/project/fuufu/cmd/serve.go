package cmd

import (
	"net/http"
	"time"

	"github.com/AugustineAurelius/fuufu/api"
	"github.com/AugustineAurelius/fuufu/config"
	"github.com/AugustineAurelius/fuufu/pkg/common"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/AugustineAurelius/fuufu/pkg/middleware"
	"github.com/AugustineAurelius/fuufu/pkg/migration"
	"github.com/AugustineAurelius/fuufu/repository"
	"github.com/AugustineAurelius/fuufu/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.0"

	"github.com/spf13/cobra"
)

func createServeCMD(manager *config.Manager) *cobra.Command {
	var debug, dev bool
	serviceName := semconv.ServiceNameKey.String("FUUFU")
	name := "fuufu"
	serveCMD := &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewWithManager(manager)

			postgresMasterConfig := manager.LoadPostgres()
			if err := migration.CheckMigrations(cmd.Context(), postgresMasterConfig); err != nil {
				return err
			}

			log.Info("migrations checked")

			pgMaster, err := common.NewPostgres(cmd.Context(), postgresMasterConfig, log)
			if err != nil {
				return err
			}
			if err = pgMaster.Pool.Ping(cmd.Context()); err != nil {
				return err
			}
			log.Info("successfully conected to postgresMaster")

			conn, err := initCollector(manager)
			if err != nil {

				return err
			}

			res, err := resource.New(cmd.Context(), resource.WithAttributes(serviceName))
			if err != nil {
				return err
			}

			shutdownTracerProvider, err := initTracerProvider(cmd.Context(), res, conn)
			if err != nil {
				return err
			}

			shutdownMetricProvider, err := initMeterProvider(cmd.Context(), res, conn)
			if err != nil {
				return err
			}
			defer func() {
				shutdownTracerProvider(cmd.Context())
				shutdownMetricProvider(cmd.Context())
			}()

			exporter, err := prometheus.New()
			if err != nil {
				panic(err)
			}

			meterProvider := metric.NewMeterProvider(
				metric.WithReader(exporter),
			)
			otel.SetMeterProvider(meterProvider)


			tracer := otel.Tracer(name)


			err = runtime.Start()
			if err != nil {
				return err
			}

			todoRepository := repository.NewCommand(&pgMaster)

			handler := server.NewHandlerMiddleware(&server.Handler{Repo: todoRepository},
				server.WithHandlerLogging(log),
				server.WithHandlerOtelMetrics(server.RegisterHandlerOtelMetrics(meterProvider)),
				server.WithhandlerTracing(tracer),
				server.WithhandlerTimeout(time.Second*5),
			)

			todoHandlers := api.NewStrictHandler(handler, nil)

			r := http.NewServeMux()
			r.Handle("/metrics", promhttp.Handler())
			h := api.HandlerFromMux(todoHandlers, r)
			h = middleware.LoggingMiddleware(log, h)
			s := &http.Server{
				Handler: h,
				Addr:    manager.LoadServer().Addr,
			}

			return s.ListenAndServe()
		},
	}

	serveCMD.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug endpoints")
	serveCMD.PersistentFlags().BoolVar(&dev, "dev", false, "Enable developer options")

	return serveCMD
}
