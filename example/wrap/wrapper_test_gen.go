package wrap

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TestInterface interface {
	Test1(a int, b float64) (param0 int, param1 error)
	Test2(a int, b float64) (param0 error)
	Test3(ctx context.Context, a int, b float64) (param0 error)
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
	m.logger.Info("call Test1", zap.Int("a", a), zap.Float64("b", b))
	defer func() { m.logger.Info("method Test1 call done", zap.Int("param0", param0), zap.Error(param1)) }()

	return m.next.Test1(a, b)
}

func (m *testLoggingMiddleware) Test2(a int, b float64) (param0 error) {
	m.logger.Info("call Test2", zap.Int("a", a), zap.Float64("b", b))
	defer func() { m.logger.Info("method Test2 call done", zap.Error(param0)) }()

	return m.next.Test2(a, b)
}

func (m *testLoggingMiddleware) Test3(ctx context.Context, a int, b float64) (param0 error) {
	m.logger.Info("call Test3", zap.Int("a", a), zap.Float64("b", b))
	defer func() { m.logger.Info("method Test3 call done", zap.Error(param0)) }()

	return m.next.Test3(ctx, a, b)
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

// Timeout
type testTimeoutMiddleware struct {
	TestInterface
	duration time.Duration
}

func WithtestTimeout(duration time.Duration) TestOption {
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
