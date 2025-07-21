package generics

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestParellelMapWithContext(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	results, err := ParallelMapWithContext(context.Background(), items, func(ctx context.Context, item int) (int, error) {
		return item * 2, nil
	})

	if err != nil {
		t.Fatalf("ParallelMapWithContext failed: %v", err)
	}

	if len(results) != len(items) {
		t.Fatalf("ParallelMapWithContext returned %d results, expected %d", len(results), len(items))
	}

}

func TestParellelMapWithError(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	expectedErr := errors.New("test error")
	results, err := ParallelMapWithContext(context.Background(), items, func(ctx context.Context, item int) (int, error) {
		return item * 2, expectedErr
	})

	if err == nil {
		t.Fatalf("ParallelMapWithContext failed: %v", err)
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("ParallelMapWithContext failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("ParallelMapWithContext returned %d results, expected 0", len(results))
	}

}

func TestParellelMapWithContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	items := []int{1, 2, 3, 4, 5}
	results, err := ParallelMapWithContext(ctx, items, func(ctx context.Context, item int) (int, error) {
		time.Sleep(10 * time.Millisecond)
		return item * 2, nil
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("ParallelMapWithContext failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("ParallelMapWithContext returned %d results, expected 0", len(results))
	}

}
