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
		contextLogging = flag.Bool("context-logging", false, "generate context logging middleware")
		includePrivate = flag.Bool("include-private", true, "include private methods in generation")
	)

	flag.Parse()

	if *name == "" {
		fmt.Println("Error: struct name is required (-n flag)")
		flag.Usage()
		os.Exit(1)
	}

	data := wrapper.StructData{
		Name:                  *name,
		MiddlewareTemplates:   make(map[string]bool, 8),
		IncludePrivateMethods: *includePrivate,
	}

	if *logging {
		data.MiddlewareTemplates["logging"] = true
		data.Logging = true
	}
	if *tracing {
		data.MiddlewareTemplates["tracing"] = true
		data.Tracing = true
	}
	if *newrelic {
		data.MiddlewareTemplates["newrelic"] = true
		data.NewRelic = true
	}
	if *timeout {
		data.MiddlewareTemplates["timeout"] = true
		data.Timeout = true
	}
	if *otelMetrics {
		data.MiddlewareTemplates["otel_metrics"] = true
		data.OtelMetrics = true
	}
	if *prometheus {
		data.MiddlewareTemplates["prometheus"] = true
		data.Prometheus = true
	}
	if *retry {
		data.MiddlewareTemplates["retry"] = true
		data.Retry = true
	}
	if *circuitBreaker {
		data.MiddlewareTemplates["circuit_breaker"] = true
		data.CircuitBreaker = true
	}
	if *contextLogging {
		data.MiddlewareTemplates["context_logging"] = true
		data.ContextLogging = true
	}

	err := wrapper.Generate(data)

	if err != nil {
		fmt.Printf("Error generating wrapper: %v\n", err)
		os.Exit(1)
	}
}
