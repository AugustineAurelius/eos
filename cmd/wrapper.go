package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/AugustineAurelius/eos/generator/wrapper"
)

func HandleWrapper() {
	var (
		name           = flag.String("name", "", "name of the struct")
		logging        = flag.Bool("logging", false, "generate logging middleware")
		tracing        = flag.Bool("tracing", false, "generate tracing middleware")
		newrelic       = flag.Bool("newrelic", false, "generate NewRelic middleware")
		timeout        = flag.Bool("timeout", false, "generate timeout middleware")
		otelMetrics    = flag.Bool("otel-metrics", false, "generate OpenTelemetry metrics middleware")
		prometheus     = flag.Bool("prometheus", false, "generate Prometheus metrics middleware")
		retry          = flag.Bool("retry", false, "generate retry middleware")
		circuitBreaker = flag.Bool("circuit-breaker", false, "generate circuit breaker middleware")
	)

	flag.Parse()

	if *name == "" {
		fmt.Println("Error: struct name is required (-n flag)")
		flag.Usage()
		os.Exit(1)
	}

	err := wrapper.Generate(wrapper.StructData{
		Name:           *name,
		Logging:        *logging,
		Tracing:        *tracing,
		NewRelic:       *newrelic,
		Timeout:        *timeout,
		OtelMetrics:    *otelMetrics,
		Prometheus:     *prometheus,
		Retry:          *retry,
		CircuitBreaker: *circuitBreaker,
	})

	if err != nil {
		fmt.Printf("Error generating wrapper: %v\n", err)
		os.Exit(1)
	}
}
