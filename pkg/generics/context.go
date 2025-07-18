package generics

import (
	"context"
)

type ContextKey[T any] struct {
	name string
}

func NewContextKey[T any](name string) ContextKey[T] {
	return ContextKey[T]{name: name}
}

func AddWithKey[T any](ctx context.Context, key ContextKey[T], value T) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetWithKey[T any](ctx context.Context, key ContextKey[T]) (T, bool) {
	value, ok := ctx.Value(key).(T)
	return value, ok
}

func GetOrDefaultWithKey[T any](ctx context.Context, key ContextKey[T], defaultValue T) T {
	if value, ok := GetWithKey(ctx, key); ok {
		return value
	}
	return defaultValue
}

type contextKey struct{}

func Add[T any](ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, contextKey{}, value)
}

func Get[T any](ctx context.Context) (T, bool) {
	value, ok := ctx.Value(contextKey{}).(T)
	return value, ok
}

func GetOrDefault[T any](ctx context.Context, defaultValue T) T {
	value, ok := Get[T](ctx)
	if !ok {
		return defaultValue
	}
	return value
}
