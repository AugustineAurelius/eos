package generics

import (
	"context"
	"sync"
)

func ParallelMap[T, R any](items []T, fn func(T) R) []R {
	return ParallelMapWithContext(context.Background(), items, fn)
}

func ParallelMapWithContext[T, R any](ctx context.Context, items []T, fn func(T) R) []R {
	if len(items) == 0 {
		return nil
	}

	results := make([]R, len(items))
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(index int, value T) {
			defer wg.Done()

			resultChan := make(chan R, 1)
			go func() {
				defer close(resultChan)
				resultChan <- fn(value)
			}()
			select {
			case <-ctx.Done():
				return
			case res := <-resultChan:
				results[index] = res
			}
		}(i, item)
	}

	wg.Wait()
	return results
}

func ParallelMapWithError[T, R any](items []T, fn func(T) (R, error)) ([]R, error) {
	return ParallelMapWithErrorAndContext(context.Background(), items, fn)
}

func ParallelMapWithErrorAndContext[T, R any](ctx context.Context, items []T, fn func(T) (R, error)) ([]R, error) {
	if len(items) == 0 {
		return nil, nil
	}

	results := make([]R, len(items))
	errors := make([]error, len(items))
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(index int, value T) {
			defer wg.Done()
			resultChan := make(chan R, 1)
			errChan := make(chan error, 1)
			go func() {
				defer close(resultChan)
				defer close(errChan)
				res, err := fn(value)
				if err != nil {
					errChan <- err
				} else {
					resultChan <- res
				}
			}()
			select {
			case <-ctx.Done():
				errors[index] = ctx.Err()
			case res := <-resultChan:
				results[index] = res
			case err := <-errChan:
				errors[index] = err
			}
		}(i, item)
	}

	wg.Wait()

	for _, err := range errors {
		if err != nil {
			return results, err
		}
	}

	return results, nil
}
