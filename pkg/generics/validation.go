package generics

import (
	"fmt"
	"reflect"
	"strings"
)

type Validator[T any] func(T) error

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

func ValidateAll[T any](value T, validators ...Validator[T]) error {
	for _, validator := range validators {
		if err := validator(value); err != nil {
			return err
		}
	}
	return nil
}

func ValidateSlice[T any](slice []T, validator Validator[T]) error {
	for i, item := range slice {
		if err := validator(item); err != nil {
			return fmt.Errorf("item[%d]: %w", i, err)
		}
	}
	return nil
}

func ValidateMap[K comparable, V any](m map[K]V, validator Validator[V]) error {
	for key, value := range m {
		if err := validator(value); err != nil {
			return fmt.Errorf("key[%v]: %w", key, err)
		}
	}
	return nil
}

func NotEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{Message: "cannot be empty"}
	}
	return nil
}

func MinLength(min int) Validator[string] {
	return func(value string) error {
		if len(value) < min {
			return &ValidationError{Message: fmt.Sprintf("must be at least %d characters long", min)}
		}
		return nil
	}
}

func MaxLength(max int) Validator[string] {
	return func(value string) error {
		if len(value) > max {
			return &ValidationError{Message: fmt.Sprintf("must be at most %d characters long", max)}
		}
		return nil
	}
}

func Length(length int) Validator[string] {
	return func(value string) error {
		if len(value) != length {
			return &ValidationError{Message: fmt.Sprintf("must be exactly %d characters long", length)}
		}
		return nil
	}
}

func Min[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](min T) Validator[T] {
	return func(value T) error {
		if value < min {
			return &ValidationError{Message: fmt.Sprintf("must be at least %v", min)}
		}
		return nil
	}
}

func Max[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](max T) Validator[T] {
	return func(value T) error {
		if value > max {
			return &ValidationError{Message: fmt.Sprintf("must be at most %v", max)}
		}
		return nil
	}
}

func Range[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](min, max T) Validator[T] {
	return func(value T) error {
		if value < min || value > max {
			return &ValidationError{Message: fmt.Sprintf("must be between %v and %v", min, max)}
		}
		return nil
	}
}

func NotNil[T any](value *T) error {
	if value == nil {
		return &ValidationError{Message: "cannot be nil"}
	}
	return nil
}

func NotZero[T comparable](value T) error {
	var zero T
	if value == zero {
		return &ValidationError{Message: "cannot be zero value"}
	}
	return nil
}

func In[T comparable](allowed ...T) Validator[T] {
	return func(value T) error {
		for _, allowedValue := range allowed {
			if value == allowedValue {
				return nil
			}
		}
		return &ValidationError{Message: fmt.Sprintf("must be one of %v", allowed)}
	}
}

func NotIn[T comparable](forbidden ...T) Validator[T] {
	return func(value T) error {
		for _, forbiddenValue := range forbidden {
			if value == forbiddenValue {
				return &ValidationError{Message: fmt.Sprintf("cannot be %v", forbiddenValue)}
			}
		}
		return nil
	}
}

func Positive[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64]() Validator[T] {
	return Min(T(0))
}

func Negative[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64]() Validator[T] {
	return func(value T) error {
		if value >= 0 {
			return &ValidationError{Message: "must be negative"}
		}
		return nil
	}
}

func SliceMinLength[T any](min int) Validator[[]T] {
	return func(value []T) error {
		if len(value) < min {
			return &ValidationError{Message: fmt.Sprintf("must have at least %d elements", min)}
		}
		return nil
	}
}

func SliceMaxLength[T any](max int) Validator[[]T] {
	return func(value []T) error {
		if len(value) > max {
			return &ValidationError{Message: fmt.Sprintf("must have at most %d elements", max)}
		}
		return nil
	}
}

func MapMinSize[K comparable, V any](min int) Validator[map[K]V] {
	return func(value map[K]V) error {
		if len(value) < min {
			return &ValidationError{Message: fmt.Sprintf("must have at least %d elements", min)}
		}
		return nil
	}
}

func MapMaxSize[K comparable, V any](max int) Validator[map[K]V] {
	return func(value map[K]V) error {
		if len(value) > max {
			return &ValidationError{Message: fmt.Sprintf("must have at most %d elements", max)}
		}
		return nil
	}
}

func ValidateStruct[T any](value T, fieldValidators map[string]Validator[T]) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return &ValidationError{Message: "value must be a struct"}
	}

	t := v.Type()
	for fieldName, validator := range fieldValidators {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return &ValidationError{Field: fieldName, Message: "field not found"}
		}

		fieldValue := reflect.New(t).Elem()
		fieldValue.Set(v)

		if err := validator(fieldValue.Interface().(T)); err != nil {
			if validationErr, ok := err.(*ValidationError); ok {
				validationErr.Field = fieldName
				return validationErr
			}
			return &ValidationError{Field: fieldName, Message: err.Error()}
		}
	}

	return nil
}
