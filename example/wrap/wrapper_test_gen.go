package wrap

import (
	"context"
	"errors"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TestInterface interface {
	Test1(a int, b float64) (param0 int, param1 error)
	Test2(a int, b float64) (param0 error)
	Test3(ctx context.Context, a int, b float64) (param0 error)
	Test4(ctx context.Context, a int, b float64) (param0 error)
}

type testCore struct {
	impl *Test
}

func (c *testCore) Test1(a int, b float64) (param0 int, param1 error) {
	return c.impl.Test1(a, b)
}

func (c *testCore) Test2(a int, b float64) (param0 error) {
	return c.impl.Test2(a, b)
}

func (c *testCore) Test3(ctx context.Context, a int, b float64) (param0 error) {
	return c.impl.Test3(ctx, a, b)
}

func (c *testCore) Test4(ctx context.Context, a int, b float64) (param0 error) {
	return c.impl.Test4(ctx, a, b)
}

// Main constructor
func NewTestMiddleware(impl *Test, opts ...TestOption) TestInterface {
	chain := TestInterface(&testCore{impl})
	for _, opt := range opts {
		chain = opt(chain)
	}
	return chain
}

// Option
type TestOption func(TestInterface) TestInterface

// Logging
type testLoggingMiddleware struct {
	next   TestInterface
	logger *zap.Logger
}

func WithTestLogging(logger *zap.Logger) TestOption {
	return func(next TestInterface) TestInterface {
		return &testLoggingMiddleware{
			next:   next,
			logger: logger.With(zap.String("struct", "Test")),
		}
	}
}

func (m *testLoggingMiddleware) Test1(a int, b float64) (param0 int, param1 error) {
	start := time.Now()
	m.logger.Info("call Test1", zap.Int("a", a), zap.Float64("b", b))
	defer func() {
		m.logger.Info("method Test1 call done", zap.Duration("diration", time.Since(start)), zap.Int("param0", param0), zap.Error(param1))
	}()

	return m.next.Test1(a, b)
}

func (m *testLoggingMiddleware) Test2(a int, b float64) (param0 error) {
	start := time.Now()
	m.logger.Info("call Test2", zap.Int("a", a), zap.Float64("b", b))
	defer func() {
		m.logger.Info("method Test2 call done", zap.Duration("diration", time.Since(start)), zap.Error(param0))
	}()

	return m.next.Test2(a, b)
}

func (m *testLoggingMiddleware) Test3(ctx context.Context, a int, b float64) (param0 error) {
	start := time.Now()
	m.logger.Info("call Test3", zap.Int("a", a), zap.Float64("b", b))
	defer func() {
		m.logger.Info("method Test3 call done", zap.Duration("diration", time.Since(start)), zap.Error(param0))
	}()

	return m.next.Test3(ctx, a, b)
}

func (m *testLoggingMiddleware) Test4(ctx context.Context, a int, b float64) (param0 error) {
	start := time.Now()
	m.logger.Info("call Test4", zap.Int("a", a), zap.Float64("b", b))
	defer func() {
		m.logger.Info("method Test4 call done", zap.Duration("diration", time.Since(start)), zap.Error(param0))
	}()

	return m.next.Test4(ctx, a, b)
}

// Tracing
type testTracingMiddleware struct {
	next   TestInterface
	tracer trace.Tracer
}

func WithTestTracing(tracer trace.Tracer) TestOption {
	return func(next TestInterface) TestInterface {
		return &testTracingMiddleware{
			next:   next,
			tracer: tracer,
		}
	}
}

func (m *testTracingMiddleware) Test1(a int, b float64) (param0 int, param1 error) {
	_, span := m.tracer.Start(context.Background(), "Test.Test1")
	defer span.End()
	return m.next.Test1(a, b)
}

func (m *testTracingMiddleware) Test2(a int, b float64) (param0 error) {
	_, span := m.tracer.Start(context.Background(), "Test.Test2")
	defer span.End()
	return m.next.Test2(a, b)
}

func (m *testTracingMiddleware) Test3(ctx context.Context, a int, b float64) (param0 error) {
	ctx, span := m.tracer.Start(ctx, "Test.Test3")
	defer span.End()
	return m.next.Test3(ctx, a, b)
}

func (m *testTracingMiddleware) Test4(ctx context.Context, a int, b float64) (param0 error) {
	ctx, span := m.tracer.Start(ctx, "Test.Test4")
	defer span.End()
	return m.next.Test4(ctx, a, b)
}

type testNewRelicTracingMiddleware struct {
	next        TestInterface
	newRelicApp *newrelic.Application
	baseLogger  *zap.Logger
}

func WithTestNewRelicTracing(app *newrelic.Application, baseLogger *zap.Logger) TestOption {
	return func(next TestInterface) TestInterface {
		return &testNewRelicTracingMiddleware{
			next:        next,
			newRelicApp: app,
			baseLogger:  baseLogger,
		}
	}
}

func (m *testNewRelicTracingMiddleware) Test1(a int, b float64) (param0 int, param1 error) {
	var logger *zap.Logger
	txn := m.newRelicApp.StartTransaction("Test.Test1")
	defer txn.End()

	bgCore, err := nrzap.WrapBackgroundCore(m.baseLogger.Core(), m.newRelicApp)
	if err == nil {
		logger = zap.New(bgCore).With(
			zap.String("method", "Test.Test1"),
			zap.String("transactionType", "background"),
		)
	} else {
		logger = m.baseLogger.With(zap.Error(err))
	}
	logger.Info("Method  Test.Test1 started")
	defer logger.Info("Method Test.Test1 completed")

	return m.next.Test1(a, b)
}

func (m *testNewRelicTracingMiddleware) Test2(a int, b float64) (param0 error) {
	var logger *zap.Logger
	txn := m.newRelicApp.StartTransaction("Test.Test2")
	defer txn.End()

	bgCore, err := nrzap.WrapBackgroundCore(m.baseLogger.Core(), m.newRelicApp)
	if err == nil {
		logger = zap.New(bgCore).With(
			zap.String("method", "Test.Test2"),
			zap.String("transactionType", "background"),
		)
	} else {
		logger = m.baseLogger.With(zap.Error(err))
	}
	logger.Info("Method  Test.Test2 started")
	defer logger.Info("Method Test.Test2 completed")

	return m.next.Test2(a, b)
}

func (m *testNewRelicTracingMiddleware) Test3(ctx context.Context, a int, b float64) (param0 error) {
	var logger *zap.Logger
	txn := newrelic.FromContext(ctx)

	if txn != nil {
		txnCore, err := nrzap.WrapTransactionCore(m.baseLogger.Core(), txn)
		if err == nil {
			logger = zap.New(txnCore).With(zap.String("method", "Test.Test3"))
		} else {
			logger = m.baseLogger.With(zap.Error(err))
		}
		seg := txn.StartSegment("Test.Test3")
		defer seg.End()
	} else {
		txn = m.newRelicApp.StartTransaction("Test.Test3")
		defer txn.End()
		ctx = newrelic.NewContext(ctx, txn)

		bgCore, err := nrzap.WrapBackgroundCore(m.baseLogger.Core(), m.newRelicApp)
		if err == nil {
			logger = zap.New(bgCore).With(
				zap.String("method", "Test.Test3"),
				zap.String("transactionType", "background"),
			)
		} else {
			logger = m.baseLogger.With(zap.Error(err))
		}
	}
	logger.Info("Method  Test.Test3 started")
	defer logger.Info("Method Test.Test3 completed")

	return m.next.Test3(ctx, a, b)
}

func (m *testNewRelicTracingMiddleware) Test4(ctx context.Context, a int, b float64) (param0 error) {
	var logger *zap.Logger
	txn := newrelic.FromContext(ctx)

	if txn != nil {
		txnCore, err := nrzap.WrapTransactionCore(m.baseLogger.Core(), txn)
		if err == nil {
			logger = zap.New(txnCore).With(zap.String("method", "Test.Test4"))
		} else {
			logger = m.baseLogger.With(zap.Error(err))
		}
		seg := txn.StartSegment("Test.Test4")
		defer seg.End()
	} else {
		txn = m.newRelicApp.StartTransaction("Test.Test4")
		defer txn.End()
		ctx = newrelic.NewContext(ctx, txn)

		bgCore, err := nrzap.WrapBackgroundCore(m.baseLogger.Core(), m.newRelicApp)
		if err == nil {
			logger = zap.New(bgCore).With(
				zap.String("method", "Test.Test4"),
				zap.String("transactionType", "background"),
			)
		} else {
			logger = m.baseLogger.With(zap.Error(err))
		}
	}
	logger.Info("Method  Test.Test4 started")
	defer logger.Info("Method Test.Test4 completed")

	return m.next.Test4(ctx, a, b)
}

// Timeout
type testTimeoutMiddleware struct {
	TestInterface
	duration time.Duration
}

// if method has context as param timeout would be applied
func WithTestTimeout(duration time.Duration) TestOption {
	return func(next TestInterface) TestInterface {
		return &testTimeoutMiddleware{
			TestInterface: next,
			duration:      duration,
		}
	}
}

func (m *testTimeoutMiddleware) Test3(ctx context.Context, a int, b float64) (param0 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.TestInterface.Test3(ctx, a, b)
}

func (m *testTimeoutMiddleware) Test4(ctx context.Context, a int, b float64) (param0 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.TestInterface.Test4(ctx, a, b)
}

type testOtelMetricsRegister struct {
	meter    metric.Meter
	Duration metric.Float64Histogram
	Calls    metric.Int64Counter
	Errors   metric.Int64Counter
	InFlight metric.Int64UpDownCounter
}

func RegisterTestOtelMetrics(provider metric.MeterProvider) *testOtelMetricsRegister {
	meter := provider.Meter("test/metrics")

	duration, _ := meter.Float64Histogram(
		"test_method_duration_seconds",
		metric.WithDescription("Method execution time distribution"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)

	calls, _ := meter.Int64Counter(
		"test_method_calls_total",
		metric.WithDescription("Total number of method calls"),
	)

	errors, _ := meter.Int64Counter(
		"test_method_errors_total",
		metric.WithDescription("Total number of method errors"),
	)

	inflight, _ := meter.Int64UpDownCounter(
		"test_method_in_flight",
		metric.WithDescription("Current number of executing methods"),
	)

	return &testOtelMetricsRegister{
		meter:    meter,
		Duration: duration,
		Calls:    calls,
		Errors:   errors,
		InFlight: inflight,
	}
}

type testOtelMetrics struct {
	TestInterface
	metrics *testOtelMetricsRegister
}

func WithTestOtelMetrics(metrics *testOtelMetricsRegister) TestOption {
	return func(next TestInterface) TestInterface {
		return &testOtelMetrics{
			TestInterface: next,
			metrics:       metrics,
		}
	}
}

func (m *testOtelMetrics) Test3(ctx context.Context, a int, b float64) (param0 error) {
	start := time.Now()
	methodName := "Test3"
	commonAttrs := []attribute.KeyValue{
		attribute.String("method", methodName),
	}

	// Track in-flight requests
	m.metrics.InFlight.Add(ctx, 1)
	defer m.metrics.InFlight.Add(ctx, -1)

	// Increment call counter
	m.metrics.Calls.Add(ctx, 1, metric.WithAttributes(commonAttrs...))

	defer func() {
		duration := time.Since(start).Seconds()
		m.metrics.Duration.Record(ctx, duration, metric.WithAttributes(commonAttrs...))

		if param0 != nil {
			errorType := param0.Error()
			switch {
			case errors.Is(param0, context.Canceled):
				errorType = "context_canceled"
			case errors.Is(param0, context.DeadlineExceeded):
				errorType = "timeout"
			}

			errorAttrs := append(commonAttrs, attribute.String("error_type", errorType))
			m.metrics.Errors.Add(ctx, 1, metric.WithAttributes(errorAttrs...))
		}

		if r := recover(); r != nil {
			errorAttrs := append(commonAttrs, attribute.String("error_type", "panic"))
			m.metrics.Errors.Add(ctx, 1, metric.WithAttributes(errorAttrs...))
			panic(r) // Re-throw panic after recording
		}
	}()

	return m.TestInterface.Test3(ctx, a, b)
}

func (m *testOtelMetrics) Test4(ctx context.Context, a int, b float64) (param0 error) {
	start := time.Now()
	methodName := "Test4"
	commonAttrs := []attribute.KeyValue{
		attribute.String("method", methodName),
	}

	// Track in-flight requests
	m.metrics.InFlight.Add(ctx, 1)
	defer m.metrics.InFlight.Add(ctx, -1)

	// Increment call counter
	m.metrics.Calls.Add(ctx, 1, metric.WithAttributes(commonAttrs...))

	defer func() {
		duration := time.Since(start).Seconds()
		m.metrics.Duration.Record(ctx, duration, metric.WithAttributes(commonAttrs...))

		if param0 != nil {
			errorType := param0.Error()
			switch {
			case errors.Is(param0, context.Canceled):
				errorType = "context_canceled"
			case errors.Is(param0, context.DeadlineExceeded):
				errorType = "timeout"
			}

			errorAttrs := append(commonAttrs, attribute.String("error_type", errorType))
			m.metrics.Errors.Add(ctx, 1, metric.WithAttributes(errorAttrs...))
		}

		if r := recover(); r != nil {
			errorAttrs := append(commonAttrs, attribute.String("error_type", "panic"))
			m.metrics.Errors.Add(ctx, 1, metric.WithAttributes(errorAttrs...))
			panic(r) // Re-throw panic after recording
		}
	}()

	return m.TestInterface.Test4(ctx, a, b)
}

type testMetrics struct {
	Duration *prometheus.HistogramVec
	Calls    *prometheus.CounterVec
	Errors   *prometheus.CounterVec
	InFlight prometheus.Gauge
}

func RegisterTestMetrics(registry prometheus.Registerer) *testMetrics {
	metrics := &testMetrics{
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "test_method_duration_seconds",
			Help:    "Method execution time distribution",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}, []string{"method"}),

		Calls: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "test_method_calls_total",
			Help: "Total number of method calls",
		}, []string{"method"}),

		Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "test_method_errors_total",
			Help: "Total number of method errors",
		}, []string{"method", "error_type"}),

		InFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "test_method_in_flight",
			Help: "Current number of executing methods",
		}),
	}

	registry.MustRegister(
		metrics.Duration,
		metrics.Calls,
		metrics.Errors,
		metrics.InFlight,
	)

	return metrics
}

type testMetricsMiddleware struct {
	next    TestInterface
	metrics *testMetrics
}

func WithTestMetrics(metrics *testMetrics) TestOption {
	return func(next TestInterface) TestInterface {
		return &testMetricsMiddleware{
			next:    next,
			metrics: metrics,
		}
	}
}

func (m *testMetricsMiddleware) Test1(a int, b float64) (param0 int, param1 error) {
	start := time.Now()
	methodName := "Test1"

	m.metrics.InFlight.Inc()
	defer m.metrics.InFlight.Dec()
	m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func() {
		duration := time.Since(start).Seconds()
		m.metrics.Duration.WithLabelValues(methodName).Observe(duration)
		if param1 != nil {
			errorType := param1.Error()
			switch {
			case errors.Is(param1, context.Canceled):
				errorType = "context_canceled"
			case errors.Is(param1, context.DeadlineExceeded):
				errorType = "timeout"
			}
			m.metrics.Errors.WithLabelValues(methodName, errorType).Inc()
		}
		if r := recover(); r != nil {
			m.metrics.Errors.WithLabelValues(methodName, "panic").Inc()
		}
	}()
	return m.next.Test1(a, b)
}

func (m *testMetricsMiddleware) Test2(a int, b float64) (param0 error) {
	start := time.Now()
	methodName := "Test2"

	m.metrics.InFlight.Inc()
	defer m.metrics.InFlight.Dec()
	m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func() {
		duration := time.Since(start).Seconds()
		m.metrics.Duration.WithLabelValues(methodName).Observe(duration)
		if param0 != nil {
			errorType := param0.Error()
			switch {
			case errors.Is(param0, context.Canceled):
				errorType = "context_canceled"
			case errors.Is(param0, context.DeadlineExceeded):
				errorType = "timeout"
			}
			m.metrics.Errors.WithLabelValues(methodName, errorType).Inc()
		}
		if r := recover(); r != nil {
			m.metrics.Errors.WithLabelValues(methodName, "panic").Inc()
		}
	}()
	return m.next.Test2(a, b)
}

func (m *testMetricsMiddleware) Test3(ctx context.Context, a int, b float64) (param0 error) {
	start := time.Now()
	methodName := "Test3"

	m.metrics.InFlight.Inc()
	defer m.metrics.InFlight.Dec()
	m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func() {
		duration := time.Since(start).Seconds()
		m.metrics.Duration.WithLabelValues(methodName).Observe(duration)
		if param0 != nil {
			errorType := param0.Error()
			switch {
			case errors.Is(param0, context.Canceled):
				errorType = "context_canceled"
			case errors.Is(param0, context.DeadlineExceeded):
				errorType = "timeout"
			}
			m.metrics.Errors.WithLabelValues(methodName, errorType).Inc()
		}
		if r := recover(); r != nil {
			m.metrics.Errors.WithLabelValues(methodName, "panic").Inc()
		}
	}()
	return m.next.Test3(ctx, a, b)
}

func (m *testMetricsMiddleware) Test4(ctx context.Context, a int, b float64) (param0 error) {
	start := time.Now()
	methodName := "Test4"

	m.metrics.InFlight.Inc()
	defer m.metrics.InFlight.Dec()
	m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func() {
		duration := time.Since(start).Seconds()
		m.metrics.Duration.WithLabelValues(methodName).Observe(duration)
		if param0 != nil {
			errorType := param0.Error()
			switch {
			case errors.Is(param0, context.Canceled):
				errorType = "context_canceled"
			case errors.Is(param0, context.DeadlineExceeded):
				errorType = "timeout"
			}
			m.metrics.Errors.WithLabelValues(methodName, errorType).Inc()
		}
		if r := recover(); r != nil {
			m.metrics.Errors.WithLabelValues(methodName, "panic").Inc()
		}
	}()
	return m.next.Test4(ctx, a, b)
}
