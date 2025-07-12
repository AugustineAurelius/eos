package generator

import (
	"github.com/AugustineAurelius/eos/generator/wrapper"
	"github.com/spf13/cobra"
)

func wrapperCMD() *cobra.Command {

	var name string
	var logging, tracing, newrelic, timeout, otelMetrics, prometheus, retry, circuitBreaker bool

	cmd := &cobra.Command{
		Use:   "wrapper",
		Short: "wrapper generator",
		Long: `Generate middleware wrappers for Go structs.
		
Available middlewares:
- logging: Add logging with zap logger
- tracing: Add OpenTelemetry tracing
- newrelic: Add NewRelic tracing and logging
- timeout: Add timeout for context-aware methods
- otel-metrics: Add OpenTelemetry metrics
- prometheus: Add Prometheus metrics
- retry: Add retry logic with exponential backoff
- circuit-breaker: Add circuit breaker pattern

If no middleware is specified, all middlewares will be generated in a single file.
If specific middlewares are selected, they will be generated in separate files.`,

		RunE: func(cmd *cobra.Command, args []string) error {

			return wrapper.Generate(wrapper.StructData{
				Name:           name,
				Logging:        logging,
				Tracing:        tracing,
				NewRelic:       newrelic,
				Timeout:        timeout,
				OtelMetrics:    otelMetrics,
				Prometheus:     prometheus,
				Retry:          retry,
				CircuitBreaker: circuitBreaker,
			})
		},
	}

	cmd.PersistentFlags().StringVarP(&name, "name", "n", "", "name of the struct")
	cmd.PersistentFlags().BoolVar(&logging, "logging", false, "generate logging middleware")
	cmd.PersistentFlags().BoolVar(&tracing, "tracing", false, "generate tracing middleware")
	cmd.PersistentFlags().BoolVar(&newrelic, "newrelic", false, "generate NewRelic middleware")
	cmd.PersistentFlags().BoolVar(&timeout, "timeout", false, "generate timeout middleware")
	cmd.PersistentFlags().BoolVar(&otelMetrics, "otel-metrics", false, "generate OpenTelemetry metrics middleware")
	cmd.PersistentFlags().BoolVar(&prometheus, "prometheus", false, "generate Prometheus metrics middleware")
	cmd.PersistentFlags().BoolVar(&retry, "retry", false, "generate retry middleware")
	cmd.PersistentFlags().BoolVar(&circuitBreaker, "circuit-breaker", false, "generate circuit breaker middleware")

	cmd.MarkPersistentFlagRequired("name")

	return cmd
}
