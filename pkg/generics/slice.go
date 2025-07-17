package generics

import (
	"cmp"
	"sort"
	"time"
)

type SliceOps[T any] []T

func FromSlice[T any](slice []T) SliceOps[T] {
	return slice
}

func (s SliceOps[T]) FilterFunc(f func(T) bool) SliceOps[T] {
	output := make(SliceOps[T], 0, len(s))
	for _, elem := range s {
		if f(elem) {
			output = append(output, elem)
		}
	}
	return output
}

func (s SliceOps[T]) FindFunc(f func(T) bool) (T, bool) {
	for _, elem := range s {
		if f(elem) {
			return elem, true
		}
	}
	var zero T
	return zero, false
}

func (s SliceOps[T]) GetFirst() (T, bool) {
	if len(s) < 1 {
		var zero T
		return zero, false
	}
	return s[0], true
}

func (s SliceOps[T]) All() Iterator[T] {
	return func(yield func(T) bool) {
		for _, elem := range s {
			if !yield(elem) {
				return
			}
		}
	}
}

func (s SliceOps[T]) SortByField(less func(T, T) bool) SliceOps[T] {
	result := make(SliceOps[T], len(s))
	copy(result, s)

	sort.Slice(result, func(i, j int) bool {
		return less(result[i], result[j])
	})

	return result
}

func (s SliceOps[T]) SortByTimeFieldAsc(fieldExtractor func(T) time.Time) SliceOps[T] {
	return s.SortByField(func(a, b T) bool {
		return fieldExtractor(a).Before(fieldExtractor(b))
	})
}

func (s SliceOps[T]) SortByTimeFieldDesc(fieldExtractor func(T) time.Time) SliceOps[T] {
	return s.SortByField(func(a, b T) bool {
		return fieldExtractor(a).After(fieldExtractor(b))
	})
}

func (s SliceOps[T]) SortByTimePtrFieldAsc(fieldExtractor func(T) *time.Time) SliceOps[T] {
	return s.SortByField(func(a, b T) bool {
		timeA := fieldExtractor(a)
		timeB := fieldExtractor(b)
		if timeA == nil && timeB == nil {
			return false
		}
		if timeA == nil {
			return true
		}
		if timeB == nil {
			return false
		}
		return timeA.Before(*timeB)
	})
}

func (s SliceOps[T]) SortByTimePtrFieldDesc(fieldExtractor func(T) *time.Time) SliceOps[T] {
	return s.SortByField(func(a, b T) bool {
		timeA := fieldExtractor(a)
		timeB := fieldExtractor(b)
		if timeA == nil && timeB == nil {
			return false
		}
		if timeA == nil {
			return false
		}
		if timeB == nil {
			return true
		}
		return timeA.After(*timeB)
	})
}

func SortByFieldAsc[T any, C cmp.Ordered](s SliceOps[T], fieldExtractor func(T) C) SliceOps[T] {
	return s.SortByField(func(a, b T) bool {
		return fieldExtractor(a) < fieldExtractor(b)
	})
}

func SortByFieldDesc[T any, C cmp.Ordered](s SliceOps[T], fieldExtractor func(T) C) SliceOps[T] {
	return s.SortByField(func(a, b T) bool {
		return fieldExtractor(a) > fieldExtractor(b)
	})
}

func FindByField[T any, K cmp.Ordered](s SliceOps[T], fieldExtractor func(T) K, value K) (T, bool) {
	for i := 0; i < len(s); i++ {
		if fieldExtractor(s[i]) == value {
			return s[i], true
		}
	}
	var zero T
	return zero, false
}

func Extract[T any, K any](s SliceOps[T], f func(T) K) SliceOps[K] {
	output := make(SliceOps[K], len(s))
	for i, elem := range s {
		output[i] = f(elem)
	}
	return output
}
