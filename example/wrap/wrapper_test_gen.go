package wrap

import (
	"context"
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)


type TestInterface interface {
    Test1(a int,b float64) (param0 int,param1 error)
    Test2(a int,b float64) (param0 error)
    Test3(ctx context.Context,a int,b float64) (param0 error)
}

type testCore struct {
	impl *Test
}


func (c *testCore) Test1(a int,b float64) (param0 int,param1 error) {
	return c.impl.Test1(a,b)
}

func (c *testCore) Test2(a int,b float64) (param0 error) {
	return c.impl.Test2(a,b)
}

func (c *testCore) Test3(ctx context.Context,a int,b float64) (param0 error) {
	return c.impl.Test3(ctx,a,b)
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


func (m *testLoggingMiddleware) Test1(a int,b float64) (param0 int,param1 error) {
    m.logger.Info("call Test1",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { m.logger.Info("method Test1 call done",zap.Int("param0", param0),zap.Error(param1),)}()

    return m.next.Test1(a,b)
}

func (m *testLoggingMiddleware) Test2(a int,b float64) (param0 error) {
    m.logger.Info("call Test2",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { m.logger.Info("method Test2 call done",zap.Error(param0),)}()

    return m.next.Test2(a,b)
}

func (m *testLoggingMiddleware) Test3(ctx context.Context,a int,b float64) (param0 error) {
    m.logger.Info("call Test3",zap.Int("a", a),zap.Float64("b", b),)
    defer func() { m.logger.Info("method Test3 call done",zap.Error(param0),)}()

    return m.next.Test3(ctx,a,b)
}


// Tracing
type testTracingMiddleware struct {
	next   TestInterface
	tracer trace.Tracer
}

func WithtestTracing(tracer trace.Tracer) TestOption {
	return func(next TestInterface) TestInterface {
		return &testTracingMiddleware{
			next:   next,
			tracer: tracer,
		}
	}
}


func (m *testTracingMiddleware)Test1 (a int,b float64) (param0 int,param1 error) {
	_, span := m.tracer.Start(context.Background(), "Test.Test1")
	defer span.End()
	return m.next.Test1(a,b)
}

func (m *testTracingMiddleware)Test2 (a int,b float64) (param0 error) {
	_, span := m.tracer.Start(context.Background(), "Test.Test2")
	defer span.End()
	return m.next.Test2(a,b)
}

func (m *testTracingMiddleware)Test3 (ctx context.Context,a int,b float64) (param0 error) {
	ctx, span := m.tracer.Start(ctx, "Test.Test3")
	defer span.End()
	return m.next.Test3(ctx,a,b)
}



// Timeout
type testTimeoutMiddleware struct {
	TestInterface
	duration time.Duration
}

func WithtestTimeout(duration time.Duration) TestOption {
	return func(next TestInterface) TestInterface {
		return &testTimeoutMiddleware{
			TestInterface:   next,
			duration: duration,
		}
	}
}




func (m *testTimeoutMiddleware)Test3 (ctx context.Context,a int,b float64) (param0 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.TestInterface.Test3(ctx,a,b)
}


type testMetrics struct {
    Duration   *prometheus.HistogramVec
    Calls      *prometheus.CounterVec
    Errors     *prometheus.CounterVec
    InFlight   prometheus.Gauge
}

func RegistertestMetrics(registry prometheus.Registerer) *testMetrics {
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


func (m *testMetricsMiddleware) Test1 (a int,b float64) (param0 int,param1 error){
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
    return m.next.Test1(a,b)
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
    return m.next.Test2(a,b)
}


func (m *testMetricsMiddleware) Test3 (ctx context.Context,a int,b float64) (param0 error){
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
    return m.next.Test3(ctx,a,b)
}

