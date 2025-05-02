package server

import (
	"context"
	"time"

	"github.com/AugustineAurelius/fuufu/api"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type HandlerInterface interface {
	GetAllTodos(ctx context.Context, request api.GetAllTodosRequestObject) (param0 api.GetAllTodosResponseObject, param1 error)
	CreateNewTask(ctx context.Context, request api.CreateNewTaskRequestObject) (param0 api.CreateNewTaskResponseObject, param1 error)
	DeleteTaskByID(ctx context.Context, request api.DeleteTaskByIDRequestObject) (param0 api.DeleteTaskByIDResponseObject, param1 error)
	GetTaskByID(ctx context.Context, request api.GetTaskByIDRequestObject) (param0 api.GetTaskByIDResponseObject, param1 error)
}

type handlerCore struct {
	impl *Handler
}

func (c *handlerCore) GetAllTodos(ctx context.Context, request api.GetAllTodosRequestObject) (param0 api.GetAllTodosResponseObject, param1 error) {
	return c.impl.GetAllTodos(ctx, request)
}

func (c *handlerCore) CreateNewTask(ctx context.Context, request api.CreateNewTaskRequestObject) (param0 api.CreateNewTaskResponseObject, param1 error) {
	return c.impl.CreateNewTask(ctx, request)
}

func (c *handlerCore) DeleteTaskByID(ctx context.Context, request api.DeleteTaskByIDRequestObject) (param0 api.DeleteTaskByIDResponseObject, param1 error) {
	return c.impl.DeleteTaskByID(ctx, request)
}

func (c *handlerCore) GetTaskByID(ctx context.Context, request api.GetTaskByIDRequestObject) (param0 api.GetTaskByIDResponseObject, param1 error) {
	return c.impl.GetTaskByID(ctx, request)
}

// Main constructor
func NewHandlerMiddleware(impl *Handler, opts ...HandlerOption) HandlerInterface {
	chain := HandlerInterface(&handlerCore{impl})
	for _, opt := range opts {
		chain = opt(chain)
	}
	return chain
}

// Option
type HandlerOption func(HandlerInterface) HandlerInterface

// Logging
type handlerLoggingMiddleware struct {
	next   HandlerInterface
	logger *zap.Logger
}

func WithHandlerLogging(logger *zap.Logger) HandlerOption {
	return func(next HandlerInterface) HandlerInterface {
		return &handlerLoggingMiddleware{
			next:   next,
			logger: logger.With(zap.String("struct", "Handler")),
		}
	}
}

func (m *handlerLoggingMiddleware) GetAllTodos(ctx context.Context, request api.GetAllTodosRequestObject) (param0 api.GetAllTodosResponseObject, param1 error) {
	m.logger.Info("call GetAllTodos")
	defer func() { m.logger.Info("method GetAllTodos call done", zap.Any("param0", param0), zap.Error(param1)) }()

	return m.next.GetAllTodos(ctx, request)
}

func (m *handlerLoggingMiddleware) CreateNewTask(ctx context.Context, request api.CreateNewTaskRequestObject) (param0 api.CreateNewTaskResponseObject, param1 error) {
	m.logger.Info("call CreateNewTask")
	defer func() { m.logger.Info("method CreateNewTask call done", zap.Any("param0", param0), zap.Error(param1)) }()

	return m.next.CreateNewTask(ctx, request)
}

func (m *handlerLoggingMiddleware) DeleteTaskByID(ctx context.Context, request api.DeleteTaskByIDRequestObject) (param0 api.DeleteTaskByIDResponseObject, param1 error) {
	m.logger.Info("call DeleteTaskByID")
	defer func() { m.logger.Info("method DeleteTaskByID call done", zap.Any("param0", param0), zap.Error(param1)) }()

	return m.next.DeleteTaskByID(ctx, request)
}

func (m *handlerLoggingMiddleware) GetTaskByID(ctx context.Context, request api.GetTaskByIDRequestObject) (param0 api.GetTaskByIDResponseObject, param1 error) {
	m.logger.Info("call GetTaskByID")
	defer func() { m.logger.Info("method GetTaskByID call done", zap.Any("param0", param0), zap.Error(param1)) }()

	return m.next.GetTaskByID(ctx, request)
}

// Tracing
type handlerTracingMiddleware struct {
	next   HandlerInterface
	tracer trace.Tracer
}

func WithhandlerTracing(tracer trace.Tracer) HandlerOption {
	return func(next HandlerInterface) HandlerInterface {
		return &handlerTracingMiddleware{
			next:   next,
			tracer: tracer,
		}
	}
}

func (m *handlerTracingMiddleware) GetAllTodos(ctx context.Context, request api.GetAllTodosRequestObject) (param0 api.GetAllTodosResponseObject, param1 error) {
	ctx, span := m.tracer.Start(ctx, "Handler.GetAllTodos")
	defer span.End()
	return m.next.GetAllTodos(ctx, request)
}

func (m *handlerTracingMiddleware) CreateNewTask(ctx context.Context, request api.CreateNewTaskRequestObject) (param0 api.CreateNewTaskResponseObject, param1 error) {
	ctx, span := m.tracer.Start(ctx, "Handler.CreateNewTask")
	defer span.End()
	return m.next.CreateNewTask(ctx, request)
}

func (m *handlerTracingMiddleware) DeleteTaskByID(ctx context.Context, request api.DeleteTaskByIDRequestObject) (param0 api.DeleteTaskByIDResponseObject, param1 error) {
	ctx, span := m.tracer.Start(ctx, "Handler.DeleteTaskByID")
	defer span.End()
	return m.next.DeleteTaskByID(ctx, request)
}

func (m *handlerTracingMiddleware) GetTaskByID(ctx context.Context, request api.GetTaskByIDRequestObject) (param0 api.GetTaskByIDResponseObject, param1 error) {
	ctx, span := m.tracer.Start(ctx, "Handler.GetTaskByID")
	defer span.End()
	return m.next.GetTaskByID(ctx, request)
}

// Timeout
type handlerTimeoutMiddleware struct {
	HandlerInterface
	duration time.Duration
}

func WithhandlerTimeout(duration time.Duration) HandlerOption {
	return func(next HandlerInterface) HandlerInterface {
		return &handlerTimeoutMiddleware{
			HandlerInterface: next,
			duration:         duration,
		}
	}
}

func (m *handlerTimeoutMiddleware) GetAllTodos(ctx context.Context, request api.GetAllTodosRequestObject) (param0 api.GetAllTodosResponseObject, param1 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.HandlerInterface.GetAllTodos(ctx, request)
}

func (m *handlerTimeoutMiddleware) CreateNewTask(ctx context.Context, request api.CreateNewTaskRequestObject) (param0 api.CreateNewTaskResponseObject, param1 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.HandlerInterface.CreateNewTask(ctx, request)
}

func (m *handlerTimeoutMiddleware) DeleteTaskByID(ctx context.Context, request api.DeleteTaskByIDRequestObject) (param0 api.DeleteTaskByIDResponseObject, param1 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.HandlerInterface.DeleteTaskByID(ctx, request)
}

func (m *handlerTimeoutMiddleware) GetTaskByID(ctx context.Context, request api.GetTaskByIDRequestObject) (param0 api.GetTaskByIDResponseObject, param1 error) {
	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()
	return m.HandlerInterface.GetTaskByID(ctx, request)
}
