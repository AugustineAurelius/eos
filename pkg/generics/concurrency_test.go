package generics

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestParallelMap(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}

	result := ParallelMap(numbers, func(n int) int {
		return n * 2
	})

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, val := range result {
		if val != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, val)
		}
	}

	empty := ParallelMap([]int{}, func(n int) int { return n * 2 })
	if len(empty) != 0 {
		t.Error("Expected empty result for empty input")
	}
}

func TestParallelMapWithContext(t *testing.T) {
	ctx := context.Background()
	numbers := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}

	result := ParallelMapWithContext(ctx, numbers, func(n int) int {
		return n * 2
	})

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, val := range result {
		if val != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, val)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result2 := ParallelMapWithContext(ctx, numbers, func(n int) int {
		time.Sleep(100 * time.Millisecond)
		return n * 2
	})

	if len(result2) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result2))
	}
}

func TestParallelMapWithError(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	expected := []string{"1", "2", "3", "4", "5"}

	result, err := ParallelMapWithError(numbers, func(n int) (string, error) {
		return string(rune(n + '0')), nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	_, err = ParallelMapWithError(numbers, func(n int) (string, error) {
		if n == 3 {
			return "", errors.New("found 3")
		}
		return string(rune(n + '0')), nil
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	result2, err := ParallelMapWithError([]int{}, func(n int) (string, error) {
		return "", nil
	})

	if err != nil {
		t.Errorf("Expected no error for empty slice, got %v", err)
	}
	if len(result2) != 0 {
		t.Error("Expected empty result for empty input")
	}
}

func TestParallelMapWithErrorAndContext(t *testing.T) {
	ctx := context.Background()
	numbers := []int{1, 2, 3, 4, 5}
	expected := []string{"1", "2", "3", "4", "5"}

	result, err := ParallelMapWithErrorAndContext(ctx, numbers, func(n int) (string, error) {
		return string(rune(n + '0')), nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	_, err = ParallelMapWithErrorAndContext(ctx, numbers, func(n int) (string, error) {
		if n == 3 {
			return "", errors.New("found 3")
		}
		return string(rune(n + '0')), nil
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result2, err := ParallelMapWithErrorAndContext(ctx, numbers, func(n int) (string, error) {
		time.Sleep(100 * time.Millisecond)
		return string(rune(n + '0')), nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result2) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result2))
	}
}

func TestConcurrencyStress(t *testing.T) {
	numbers := make([]int, 100)
	for i := range numbers {
		numbers[i] = i
	}

	result := ParallelMap(numbers, func(n int) int {
		return n * 2
	})

	if len(result) != len(numbers) {
		t.Errorf("Expected %d results, got %d", len(numbers), len(result))
	}

	for i, val := range result {
		if val != i*2 {
			t.Errorf("Expected %d at index %d, got %d", i*2, i, val)
		}
	}
}

func TestParallelMapWithContext_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result := ParallelMapWithContext(ctx, []int{1, 2, 3, 4, 5}, func(n int) int {
		time.Sleep(100 * time.Millisecond * time.Duration(n))
		return n * 2
	})

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
	}

	if !reflect.DeepEqual(result, []int{2, 0, 0, 0, 0}) && !reflect.DeepEqual(result, []int{0, 0, 0, 0, 0}) {
		t.Errorf("Expected empty result, got %v", result)
	}
}
