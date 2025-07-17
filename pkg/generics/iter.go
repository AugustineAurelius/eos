package generics

import (
	"iter"
	"slices"
)

type Iterator[T any] iter.Seq[T]

func NewIterator[T any](slice []T) Iterator[T] {
	return func(yield func(T) bool) {
		for _, elem := range slice {
			if !yield(elem) {
				return
			}
		}
	}
}

func (i Iterator[T]) FilterFunc(predicate func(T) bool) Iterator[T] {
	return func(yield func(T) bool) {
		for elem := range i {
			if predicate(elem) && !yield(elem) {
				return
			}
		}
	}
}

func (i Iterator[T]) Map(transform func(T) T) Iterator[T] {
	return func(yield func(T) bool) {
		for elem := range i {
			if !yield(transform(elem)) {
				return
			}
		}
	}
}

func (i Iterator[T]) Take(n int) Iterator[T] {
	count := 0
	return func(yield func(T) bool) {
		for elem := range i {
			if count >= n || !yield(elem) {
				return
			}
			count++
		}
	}
}

func (i Iterator[T]) Find(findFunc func(T) bool) Iterator[T] {
	return func(yield func(T) bool) {
		for elem := range i {
			if findFunc(elem) && !yield(elem) {
				return
			}
		}
	}
}

func (i Iterator[T]) Distinct(keyFunc func(T) any) Iterator[T] {
	seen := make(map[any]bool)
	return func(yield func(T) bool) {
		for elem := range i {
			key := keyFunc(elem)
			if !seen[key] {
				seen[key] = true
				if !yield(elem) {
					return
				}
			}
		}
	}
}

func (i Iterator[T]) First() (T, bool) {
	for elem := range i {
		return elem, true
	}
	var zero T
	return zero, false
}

func (i Iterator[T]) ForEach(f func(T)) {
	for elem := range i {
		f(elem)
	}
}

func (i Iterator[T]) Collect() []T {
	return slices.Collect(iter.Seq[T](i))
}

func (i Iterator[T]) Sort(sortFunc func(x, y T) int) []T {
	return slices.SortedFunc(iter.Seq[T](i), sortFunc)
}

func (i Iterator[T]) FilterByField(fieldExtractor func(T) any, value any) Iterator[T] {
	return func(yield func(T) bool) {
		for elem := range i {
			if fieldExtractor(elem) == value && !yield(elem) {
				return
			}
		}
	}
}

func (i Iterator[T]) ExtractField(fieldExtractor func(T) any) Iterator[any] {
	return func(yield func(any) bool) {
		for elem := range i {
			if !yield(fieldExtractor(elem)) {
				return
			}
		}
	}
}

func MapTo[T, R any](i Iterator[T], transform func(T) R) Iterator[R] {
	return func(yield func(R) bool) {
		for elem := range i {
			if !yield(transform(elem)) {
				return
			}
		}
	}
}

type IterWithErr[T any] iter.Seq2[T, error]
