package wrap

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"sync"
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
    Test4(ctx context.Context,a int,b float64) (param0 error)
    Test1(a int,b *Test222) (param0 int,param1 error)
    Test2(a int,b float64) (param0 error)
    Test3(ctx context.Context,c int,b float64) (param0 error)
    Test5(ctx context.Context,a int,b float64) (param0 int,param1 error)
}

type testCore struct {
	impl *Test
}


func (core *testCore) Test4(ctx context.Context,a int,b float64) (param0 error) {
	return core.impl.Test4(ctx,a,b)
}

func (core *testCore) Test1(a int,b *Test222) (param0 int,param1 error) {
	return core.impl.Test1(a,b)
}

func (core *testCore) Test2(a int,b float64) (param0 error) {
	return core.impl.Test2(a,b)
}

func (core *testCore) Test3(ctx context.Context,c int,b float64) (param0 error) {
	return core.impl.Test3(ctx,c,b)
}

func (core *testCore) Test5(ctx context.Context,a int,b float64) (param0 int,param1 error) {
	return core.impl.Test5(ctx,a,b)
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

//Logging
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


func (m *testLoggingMiddleware) Test4(ctx context.Context,a int,b float64) (param0 error) {
    start := time.Now()
    m.logger.Info("call Test4",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { m.logger.Info("method Test4 call done", zap.Duration("diration", time.Since(start)),zap.Error(param0),)}()

    return m.next.Test4(ctx,a,b)
}

func (m *testLoggingMiddleware) Test1(a int,b *Test222) (param0 int,param1 error) {
    start := time.Now()
    m.logger.Info("call Test1",zap.Int("a", a),zap.Any("b", b),)
    defer func() { m.logger.Info("method Test1 call done", zap.Duration("diration", time.Since(start)),zap.Int("param0", param0),zap.Error(param1),)}()

    return m.next.Test1(a,b)
}

func (m *testLoggingMiddleware) Test2(a int,b float64) (param0 error) {
    start := time.Now()
    m.logger.Info("call Test2",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { m.logger.Info("method Test2 call done", zap.Duration("diration", time.Since(start)),zap.Error(param0),)}()

    return m.next.Test2(a,b)
}

func (m *testLoggingMiddleware) Test3(ctx context.Context,c int,b float64) (param0 error) {
    start := time.Now()
    m.logger.Info("call Test3",zap.Int("c", c),zap.Float64("b", b),)
    defer func() { m.logger.Info("method Test3 call done", zap.Duration("diration", time.Since(start)),zap.Error(param0),)}()

    return m.next.Test3(ctx,c,b)
}

func (m *testLoggingMiddleware) Test5(ctx context.Context,a int,b float64) (param0 int,param1 error) {
    start := time.Now()
    m.logger.Info("call Test5",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { m.logger.Info("method Test5 call done", zap.Duration("diration", time.Since(start)),zap.Int("param0", param0),zap.Error(param1),)}()

    return m.next.Test5(ctx,a,b)
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


func (m *testTracingMiddleware)Test4 (ctx context.Context,a int,b float64) (param0 error) {
	ctx, span := m.tracer.Start(ctx, "Test.Test4")
	defer span.End()
	return m.next.Test4(ctx,a,b)
}

func (m *testTracingMiddleware)Test1 (a int,b *Test222) (param0 int,param1 error) {
	_, span := m.tracer.Start(context.Background(), "Test.Test1")
	defer span.End()
	return m.next.Test1(a,b)
}

func (m *testTracingMiddleware)Test2 (a int,b float64) (param0 error) {
	_, span := m.tracer.Start(context.Background(), "Test.Test2")
	defer span.End()
	return m.next.Test2(a,b)
}

func (m *testTracingMiddleware)Test3 (ctx context.Context,c int,b float64) (param0 error) {
	ctx, span := m.tracer.Start(ctx, "Test.Test3")
	defer span.End()
	return m.next.Test3(ctx,c,b)
}

func (m *testTracingMiddleware)Test5 (ctx context.Context,a int,b float64) (param0 int,param1 error) {
	ctx, span := m.tracer.Start(ctx, "Test.Test5")
	defer span.End()
	return m.next.Test5(ctx,a,b)
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


func (m *testNewRelicTracingMiddleware) Test4 (ctx context.Context,a int,b float64) (param0 error) {
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
	logger.Info("call TestTest4",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { logger.Info("method TestTest4 call done",zap.Error(param0),)}()


	return m.next.Test4(ctx,a,b)
}

func (m *testNewRelicTracingMiddleware) Test1 (a int,b *Test222) (param0 int,param1 error) {
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
	logger.Info("call TestTest1",zap.Int("a", a),zap.Any("b", b),)
    defer func() { logger.Info("method TestTest1 call done",zap.Int("param0", param0),zap.Error(param1),)}()


	return m.next.Test1(a,b)
}

func (m *testNewRelicTracingMiddleware) Test2 (a int,b float64) (param0 error) {
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
	logger.Info("call TestTest2",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { logger.Info("method TestTest2 call done",zap.Error(param0),)}()


	return m.next.Test2(a,b)
}

func (m *testNewRelicTracingMiddleware) Test3 (ctx context.Context,c int,b float64) (param0 error) {
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
	logger.Info("call TestTest3",zap.Int("c", c),zap.Float64("b", b),)
    defer func() { logger.Info("method TestTest3 call done",zap.Error(param0),)}()


	return m.next.Test3(ctx,c,b)
}

func (m *testNewRelicTracingMiddleware) Test5 (ctx context.Context,a int,b float64) (param0 int,param1 error) {
	var logger *zap.Logger
	txn := newrelic.FromContext(ctx)
	
	if txn != nil {
		txnCore, err := nrzap.WrapTransactionCore(m.baseLogger.Core(), txn)
		if err == nil {
			logger = zap.New(txnCore).With(zap.String("method", "Test.Test5"))
		} else {
			logger = m.baseLogger.With(zap.Error(err))
		}
		seg := txn.StartSegment("Test.Test5")
		defer seg.End()
	} else {
		txn = m.newRelicApp.StartTransaction("Test.Test5")
		defer txn.End()
		ctx = newrelic.NewContext(ctx, txn)
		
		bgCore, err := nrzap.WrapBackgroundCore(m.baseLogger.Core(), m.newRelicApp)
		if err == nil {
			logger = zap.New(bgCore).With(
				zap.String("method", "Test.Test5"),
				zap.String("transactionType", "background"),
			)
		} else {
			logger = m.baseLogger.With(zap.Error(err))
		}
	}
	logger.Info("call TestTest5",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { logger.Info("method TestTest5 call done",zap.Int("param0", param0),zap.Error(param1),)}()


	return m.next.Test5(ctx,a,b)
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
			TestInterface:   next,
			duration: duration,
		}
	}
}


func (m *testTimeoutMiddleware)Test4 (ctx context.Context,a int,b float64) (param0 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.TestInterface.Test4(ctx,a,b)
}



func (m *testTimeoutMiddleware)Test3 (ctx context.Context,c int,b float64) (param0 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.TestInterface.Test3(ctx,c,b)
}

func (m *testTimeoutMiddleware)Test5 (ctx context.Context,a int,b float64) (param0 int,param1 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.TestInterface.Test5(ctx,a,b)
}

type testOtelMetricsRegister struct {
    meter       metric.Meter
    Duration    metric.Float64Histogram
    Calls       metric.Int64Counter
    Errors      metric.Int64Counter
    InFlight    metric.Int64UpDownCounter
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
        meter:     meter,
        Duration:  duration,
        Calls:     calls,
        Errors:    errors,
        InFlight:  inflight,
    }
}

type testOtelMetrics struct {
    TestInterface
    metrics *testOtelMetricsRegister
}

func WithTestOtelMetrics(metrics *testOtelMetricsRegister) TestOption {
    return func(next TestInterface) TestInterface {
        return &testOtelMetrics{
            TestInterface:    next,
            metrics: metrics,
        }
    }
}


func (m *testOtelMetrics) Test4 (ctx context.Context,a int,b float64) (param0 error){
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

    return m.TestInterface.Test4(ctx,a,b)
}



func (m *testOtelMetrics) Test3 (ctx context.Context,c int,b float64) (param0 error){
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

    return m.TestInterface.Test3(ctx,c,b)
}

func (m *testOtelMetrics) Test5 (ctx context.Context,a int,b float64) (param0 int,param1 error){
    start := time.Now()
    methodName := "Test5"
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

        if param1 != nil {
            errorType := param1.Error()
            switch {
            case errors.Is(param1, context.Canceled):
                errorType = "context_canceled"
            case errors.Is(param1, context.DeadlineExceeded):
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

    return m.TestInterface.Test5(ctx,a,b)
}



type testMetrics struct {
    Duration   *prometheus.HistogramVec
    Calls      *prometheus.CounterVec
    Errors     *prometheus.CounterVec
    InFlight   prometheus.Gauge
}

func RegisterTestMetrics(registry prometheus.Registerer) *testMetrics {
    metrics := &testMetrics{
        Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
            Name: "test_method_duration_seconds",
            Help: "Method execution time distribution",
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
    TestInterface
    metrics *testMetrics
}

func WithTestMetrics(metrics *testMetrics) TestOption {
    return func(next TestInterface) TestInterface {
        return &testMetricsMiddleware{
            TestInterface:    next,
            metrics: metrics,
        }
    }
}


func (m *testMetricsMiddleware) Test4 (ctx context.Context,a int,b float64) (param0 error){
    start := time.Now()
    methodName := "Test4"
    
    m.metrics.InFlight.Inc()
    defer m.metrics.InFlight.Dec()
    m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func(){
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
    return m.TestInterface.Test4(ctx,a,b)
}


func (m *testMetricsMiddleware) Test1 (a int,b *Test222) (param0 int,param1 error){
    start := time.Now()
    methodName := "Test1"
    
    m.metrics.InFlight.Inc()
    defer m.metrics.InFlight.Dec()
    m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func(){
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
    return m.TestInterface.Test1(a,b)
}


func (m *testMetricsMiddleware) Test2 (a int,b float64) (param0 error){
    start := time.Now()
    methodName := "Test2"
    
    m.metrics.InFlight.Inc()
    defer m.metrics.InFlight.Dec()
    m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func(){
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
    return m.TestInterface.Test2(a,b)
}


func (m *testMetricsMiddleware) Test3 (ctx context.Context,c int,b float64) (param0 error){
    start := time.Now()
    methodName := "Test3"
    
    m.metrics.InFlight.Inc()
    defer m.metrics.InFlight.Dec()
    m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func(){
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
    return m.TestInterface.Test3(ctx,c,b)
}


func (m *testMetricsMiddleware) Test5 (ctx context.Context,a int,b float64) (param0 int,param1 error){
    start := time.Now()
    methodName := "Test5"
    
    m.metrics.InFlight.Inc()
    defer m.metrics.InFlight.Dec()
    m.metrics.Calls.WithLabelValues(methodName).Inc()

	defer func(){
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
    return m.TestInterface.Test5(ctx,a,b)
}





type testRetryMiddleware struct {
	TestInterface
	retryConfig *retryConfig
}

type retryConfig struct {
	timeout *time.Duration
	//default 100 ms
	delay time.Duration
	//default 5 s
	maxDelay time.Duration
	// default 5 attempts
	attempts int
	// if false -> stop retry default return always true
	shouldRetryAfter func(error) bool
}
type retryOpt func(*retryConfig)

func NewRetryConfig(opts ...retryOpt) *retryConfig {
	cfg := &retryConfig{
		timeout:nil,
		delay:             100 * time.Millisecond,
		maxDelay:          5 * time.Second,
		attempts:          5,
		shouldRetryAfter:  func(err error) bool { return true },
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func RetryWithAttempts(attempts int) retryOpt {
	return func(rc *retryConfig) {
		rc.attempts = attempts
	}
}
func RetryWithDelay(delay time.Duration) retryOpt {
	return func(rc *retryConfig) {
		rc.delay = delay
	}
}
func RetryWithMaximumDelay(maxDelay time.Duration) retryOpt {
	return func(rc *retryConfig) {
		rc.maxDelay = maxDelay
	}
}
func RetryWithTimeout(timeout time.Duration) retryOpt {
	return func(rc *retryConfig) {
		rc.timeout = &timeout
	}
}
func RetryWithShouldRetryAfter(shouldRetryAfter func(error) bool) retryOpt {
	return func(rc *retryConfig) {
		rc.shouldRetryAfter = shouldRetryAfter
	}
}


func WithTestRetry(cfg *retryConfig) TestOption  {
    return func(next TestInterface) TestInterface {
		return &testRetryMiddleware{
			TestInterface:        next,
			retryConfig: cfg,
		}
	}
}


func (m *testRetryMiddleware) Test4 (ctx context.Context,a int,b float64) (param0 error){
    var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok && m.retryConfig.timeout != nil {
		ctx, cancel = context.WithTimeout(ctx, *m.retryConfig.timeout)
		defer cancel()
	}
	for attempt := 0; attempt < m.retryConfig.attempts; attempt++ {
        if ctx.Err() != nil {
			param0 = ctx.Err()
			return 
		}
	

		param0 = m.TestInterface.Test4(ctx,a,b)
		if param0 == nil {
			return 
		}

        if !m.retryConfig.shouldRetryAfter(param0) {
			return
		}

		if attempt == m.retryConfig.attempts-1 {
			return
		}

		backoff := m.retryConfig.delay * time.Duration(math.Pow(2, float64(attempt)))
		if backoff > m.retryConfig.maxDelay {
			backoff = m.retryConfig.maxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoff)))

		time.Sleep(jitter)
	}
    return 
}


func (m *testRetryMiddleware) Test1 (a int,b *Test222) (param0 int,param1 error){
	for attempt := 0; attempt < m.retryConfig.attempts; attempt++ {
	

		param0,param1 = m.TestInterface.Test1(a,b)
		if param1 == nil {
			return 
		}

        if !m.retryConfig.shouldRetryAfter(param1) {
			return
		}

		if attempt == m.retryConfig.attempts-1 {
			return
		}

		backoff := m.retryConfig.delay * time.Duration(math.Pow(2, float64(attempt)))
		if backoff > m.retryConfig.maxDelay {
			backoff = m.retryConfig.maxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoff)))

		time.Sleep(jitter)
	}
    return 
}


func (m *testRetryMiddleware) Test2 (a int,b float64) (param0 error){
	for attempt := 0; attempt < m.retryConfig.attempts; attempt++ {
	

		param0 = m.TestInterface.Test2(a,b)
		if param0 == nil {
			return 
		}

        if !m.retryConfig.shouldRetryAfter(param0) {
			return
		}

		if attempt == m.retryConfig.attempts-1 {
			return
		}

		backoff := m.retryConfig.delay * time.Duration(math.Pow(2, float64(attempt)))
		if backoff > m.retryConfig.maxDelay {
			backoff = m.retryConfig.maxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoff)))

		time.Sleep(jitter)
	}
    return 
}


func (m *testRetryMiddleware) Test3 (ctx context.Context,c int,b float64) (param0 error){
    var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok && m.retryConfig.timeout != nil {
		ctx, cancel = context.WithTimeout(ctx, *m.retryConfig.timeout)
		defer cancel()
	}
	for attempt := 0; attempt < m.retryConfig.attempts; attempt++ {
        if ctx.Err() != nil {
			param0 = ctx.Err()
			return 
		}
	

		param0 = m.TestInterface.Test3(ctx,c,b)
		if param0 == nil {
			return 
		}

        if !m.retryConfig.shouldRetryAfter(param0) {
			return
		}

		if attempt == m.retryConfig.attempts-1 {
			return
		}

		backoff := m.retryConfig.delay * time.Duration(math.Pow(2, float64(attempt)))
		if backoff > m.retryConfig.maxDelay {
			backoff = m.retryConfig.maxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoff)))

		time.Sleep(jitter)
	}
    return 
}


func (m *testRetryMiddleware) Test5 (ctx context.Context,a int,b float64) (param0 int,param1 error){
    var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok && m.retryConfig.timeout != nil {
		ctx, cancel = context.WithTimeout(ctx, *m.retryConfig.timeout)
		defer cancel()
	}
	for attempt := 0; attempt < m.retryConfig.attempts; attempt++ {
        if ctx.Err() != nil {
			param1 = ctx.Err()
			return 
		}
	

		param0,param1 = m.TestInterface.Test5(ctx,a,b)
		if param1 == nil {
			return 
		}

        if !m.retryConfig.shouldRetryAfter(param1) {
			return
		}

		if attempt == m.retryConfig.attempts-1 {
			return
		}

		backoff := m.retryConfig.delay * time.Duration(math.Pow(2, float64(attempt)))
		if backoff > m.retryConfig.maxDelay {
			backoff = m.retryConfig.maxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoff)))

		time.Sleep(jitter)
	}
    return 
}






type circuitBreakerConfig struct {
	// default 5
	errorsAmountToOpen  int
	succesAmountToClose int
	// default 5 s
	openInterval time.Duration
	// if false -> error not counts
	shouldCountAfter func(error) bool
}
type circuitBreakerOpt func(*circuitBreakerConfig)

func NewCircuitBreakerConfig(opts ...circuitBreakerOpt) *circuitBreakerConfig {
	cfg := &circuitBreakerConfig{
		errorsAmountToOpen:  5,
		succesAmountToClose: 5,
		openInterval:        5 * time.Second,
		shouldCountAfter:    func(err error) bool { return true },
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func CircuitBreakerWithOpenInterval(openInterval time.Duration) circuitBreakerOpt {
	return func(cb *circuitBreakerConfig) {
		cb.openInterval = openInterval
	}
}
func CircuitBreakerWithSucessAmountToClose (succesAmountToClose int) circuitBreakerOpt {
	return func(cb *circuitBreakerConfig) {
		cb.succesAmountToClose = succesAmountToClose
	}
}
func CircuitBreakerWithErrorsAmountToOpen(errorsAmountToOpen int) circuitBreakerOpt {
	return func(cb *circuitBreakerConfig) {
		cb.errorsAmountToOpen = errorsAmountToOpen
	}
}
func CircuitBreakerWithShouldCountAfter(shouldCountAfter func(error) bool) circuitBreakerOpt {
	return func(cb *circuitBreakerConfig) {
		cb.shouldCountAfter = shouldCountAfter
	}
}

type circuitBreakerState int

const (
	Closed circuitBreakerState = iota
	HalfOpen
	Open
)


var ErrOpenCircuitBreaker = errors.New("TestCircuitBreaker: circuit is open")

type testCircuitBreakerMiddleware struct {
	TestInterface

	mu                     sync.Mutex
	currentAmountOfErrors  int
	currentAmountOfSuccess int
	state                  circuitBreakerState
	openedAt               *time.Time
	cfg                    *circuitBreakerConfig
}


func WithTestCircuitBreaker(cfg *circuitBreakerConfig) TestOption  {
    return func(next TestInterface) TestInterface {
		return &testCircuitBreakerMiddleware{
			TestInterface:        next,
            state: Closed,
			cfg: cfg,
		}
	}
}


func (m *testCircuitBreakerMiddleware)  Test4 (ctx context.Context,a int,b float64) (param0 error) {
    m.mu.Lock()
	defer m.mu.Unlock()
	
	switch m.state {
	case Open:
		if m.openedAt != nil && m.openedAt.After(time.Now()) {
			param0 = ErrOpenCircuitBreaker
			return
		}
		param0 = m.TestInterface.Test4(ctx,a,b)
		if param0 == nil {
			m.state = HalfOpen
			m.currentAmountOfSuccess = 1
			m.openedAt = nil
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
	case HalfOpen:
		param0 = m.TestInterface.Test4(ctx,a,b)
		if param0 == nil {
			m.currentAmountOfSuccess++
			if m.currentAmountOfSuccess >= m.cfg.succesAmountToClose {
				m.state = Closed
				m.currentAmountOfErrors = 0
			}
			return
		}
		if !m.cfg.shouldCountAfter(param0) {
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
		m.state = Open
	case Closed:
		param0 = m.TestInterface.Test4(ctx,a,b)
		if param0 == nil {
			m.currentAmountOfErrors = 0
			m.openedAt = nil
			return
		}

		if !m.cfg.shouldCountAfter(param0) {
			return
		}

		m.currentAmountOfErrors++

		if m.currentAmountOfErrors >= m.cfg.errorsAmountToOpen {
			openedAt := time.Now().Add(m.cfg.openInterval)
			m.openedAt = &openedAt
			m.state = Open
		}
	}

	return
	
}


func (m *testCircuitBreakerMiddleware)  Test1 (a int,b *Test222) (param0 int,param1 error) {
    m.mu.Lock()
	defer m.mu.Unlock()
	
	switch m.state {
	case Open:
		if m.openedAt != nil && m.openedAt.After(time.Now()) {
			param1 = ErrOpenCircuitBreaker
			return
		}
		param0,param1 = m.TestInterface.Test1(a,b)
		if param1 == nil {
			m.state = HalfOpen
			m.currentAmountOfSuccess = 1
			m.openedAt = nil
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
	case HalfOpen:
		param0,param1 = m.TestInterface.Test1(a,b)
		if param1 == nil {
			m.currentAmountOfSuccess++
			if m.currentAmountOfSuccess >= m.cfg.succesAmountToClose {
				m.state = Closed
				m.currentAmountOfErrors = 0
			}
			return
		}
		if !m.cfg.shouldCountAfter(param1) {
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
		m.state = Open
	case Closed:
		param0,param1 = m.TestInterface.Test1(a,b)
		if param1 == nil {
			m.currentAmountOfErrors = 0
			m.openedAt = nil
			return
		}

		if !m.cfg.shouldCountAfter(param1) {
			return
		}

		m.currentAmountOfErrors++

		if m.currentAmountOfErrors >= m.cfg.errorsAmountToOpen {
			openedAt := time.Now().Add(m.cfg.openInterval)
			m.openedAt = &openedAt
			m.state = Open
		}
	}

	return
	
}


func (m *testCircuitBreakerMiddleware)  Test2 (a int,b float64) (param0 error) {
    m.mu.Lock()
	defer m.mu.Unlock()
	
	switch m.state {
	case Open:
		if m.openedAt != nil && m.openedAt.After(time.Now()) {
			param0 = ErrOpenCircuitBreaker
			return
		}
		param0 = m.TestInterface.Test2(a,b)
		if param0 == nil {
			m.state = HalfOpen
			m.currentAmountOfSuccess = 1
			m.openedAt = nil
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
	case HalfOpen:
		param0 = m.TestInterface.Test2(a,b)
		if param0 == nil {
			m.currentAmountOfSuccess++
			if m.currentAmountOfSuccess >= m.cfg.succesAmountToClose {
				m.state = Closed
				m.currentAmountOfErrors = 0
			}
			return
		}
		if !m.cfg.shouldCountAfter(param0) {
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
		m.state = Open
	case Closed:
		param0 = m.TestInterface.Test2(a,b)
		if param0 == nil {
			m.currentAmountOfErrors = 0
			m.openedAt = nil
			return
		}

		if !m.cfg.shouldCountAfter(param0) {
			return
		}

		m.currentAmountOfErrors++

		if m.currentAmountOfErrors >= m.cfg.errorsAmountToOpen {
			openedAt := time.Now().Add(m.cfg.openInterval)
			m.openedAt = &openedAt
			m.state = Open
		}
	}

	return
	
}


func (m *testCircuitBreakerMiddleware)  Test3 (ctx context.Context,c int,b float64) (param0 error) {
    m.mu.Lock()
	defer m.mu.Unlock()
	
	switch m.state {
	case Open:
		if m.openedAt != nil && m.openedAt.After(time.Now()) {
			param0 = ErrOpenCircuitBreaker
			return
		}
		param0 = m.TestInterface.Test3(ctx,c,b)
		if param0 == nil {
			m.state = HalfOpen
			m.currentAmountOfSuccess = 1
			m.openedAt = nil
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
	case HalfOpen:
		param0 = m.TestInterface.Test3(ctx,c,b)
		if param0 == nil {
			m.currentAmountOfSuccess++
			if m.currentAmountOfSuccess >= m.cfg.succesAmountToClose {
				m.state = Closed
				m.currentAmountOfErrors = 0
			}
			return
		}
		if !m.cfg.shouldCountAfter(param0) {
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
		m.state = Open
	case Closed:
		param0 = m.TestInterface.Test3(ctx,c,b)
		if param0 == nil {
			m.currentAmountOfErrors = 0
			m.openedAt = nil
			return
		}

		if !m.cfg.shouldCountAfter(param0) {
			return
		}

		m.currentAmountOfErrors++

		if m.currentAmountOfErrors >= m.cfg.errorsAmountToOpen {
			openedAt := time.Now().Add(m.cfg.openInterval)
			m.openedAt = &openedAt
			m.state = Open
		}
	}

	return
	
}


func (m *testCircuitBreakerMiddleware)  Test5 (ctx context.Context,a int,b float64) (param0 int,param1 error) {
    m.mu.Lock()
	defer m.mu.Unlock()
	
	switch m.state {
	case Open:
		if m.openedAt != nil && m.openedAt.After(time.Now()) {
			param1 = ErrOpenCircuitBreaker
			return
		}
		param0,param1 = m.TestInterface.Test5(ctx,a,b)
		if param1 == nil {
			m.state = HalfOpen
			m.currentAmountOfSuccess = 1
			m.openedAt = nil
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
	case HalfOpen:
		param0,param1 = m.TestInterface.Test5(ctx,a,b)
		if param1 == nil {
			m.currentAmountOfSuccess++
			if m.currentAmountOfSuccess >= m.cfg.succesAmountToClose {
				m.state = Closed
				m.currentAmountOfErrors = 0
			}
			return
		}
		if !m.cfg.shouldCountAfter(param1) {
			return
		}
		openedAt := time.Now().Add(m.cfg.openInterval)
		m.openedAt = &openedAt
		m.state = Open
	case Closed:
		param0,param1 = m.TestInterface.Test5(ctx,a,b)
		if param1 == nil {
			m.currentAmountOfErrors = 0
			m.openedAt = nil
			return
		}

		if !m.cfg.shouldCountAfter(param1) {
			return
		}

		m.currentAmountOfErrors++

		if m.currentAmountOfErrors >= m.cfg.errorsAmountToOpen {
			openedAt := time.Now().Add(m.cfg.openInterval)
			m.openedAt = &openedAt
			m.state = Open
		}
	}

	return
	
}

