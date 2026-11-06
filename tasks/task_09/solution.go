package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func worker[T any, R any](ctx context.Context, jobs chan int, in []T, results []Result[R], fn func(context.Context, T) (R, error)) {
	for i := range jobs {
		if ctx.Err() != nil {
			return
		}
		value, err := fn(ctx, in[i])
		if err != nil {
			results[i].Err = err
		} else {
			results[i].Value = value
		}
	}
}

func sendJobs(ctx context.Context, jobs chan<- int, count int) {
	defer close(jobs)
	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			return
		case jobs <- i:
		}
	}
}

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	switch {
	case workers <= 0:
		return nil, ErrInvalidWorkers
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case len(in) == 0:
		return make([]Result[R], 0), nil
	}

	results := make([]Result[R], len(in))
	jobs := make(chan int)
	var wg sync.WaitGroup

	for w := 0; w < min(workers, len(in)); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, jobs, in, results, fn)
		}()
	}

	sendJobs(ctx, jobs, len(in))
	wg.Wait()
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
