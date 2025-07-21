package generics

import (
	"context"
	"sync"
)

func ParallelMap[T, R any](items []T, fn func(context.Context, T) (R, error)) ([]R, error) {
	return ParallelMapWithContext(context.Background(), items, fn)
}

func ParallelMapWithContext[T, R any](ctx context.Context, items []T, fn func(context.Context, T) (R, error)) ([]R, error) {
	if len(items) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	resultChan := make(chan R, len(items))
	errChan := make(chan error, len(items))
	var wg sync.WaitGroup
	wg.Add(len(items))

	for i, item := range items {
		go func(index int, value T) {
			defer wg.Done()

			doneChan := make(chan struct{}, 1)
			var result R
			var err error
			go func() {
				result, err = fn(ctx, value)
				close(doneChan)
			}()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
			case <-doneChan:
				if err != nil {
					errChan <- err
					cancel(err)
				} else {
					resultChan <- result
				}
			}
		}(i, item)
	}

	wg.Wait()
	close(resultChan)
	close(errChan)

	select {
	case err := <-errChan:
		if err != nil {
			return nil, context.Cause(ctx)
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	results := make([]R, 0, len(items))
	for result := range resultChan {
		results = append(results, result)
	}

	return results, nil
}
