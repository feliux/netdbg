package revdns

import (
	"context"
	"sync"
)

// Executor defines the interface for executing reverse DNS lookups.
type Executor interface {
	Execute(ctx context.Context, opts *Options) <-chan Result
}

// DefaultExecutor implements Executor with concurrent worker pool.
type DefaultExecutor struct{}

// Execute performs reverse DNS lookups concurrently based on the provided options.
func (e *DefaultExecutor) Execute(ctx context.Context, opts *Options) <-chan Result {
	results := make(chan Result)
	work := make(chan string)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < opts.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.worker(ctx, opts, work, results)
		}()
	}

	// Feed work
	go func() {
		defer close(work)
		if opts.File != "" {
			if err := feedFromFile(opts.File, work); err != nil {
				results <- Result{Error: err}
			}
		} else if opts.Addr != "" {
			work <- opts.Addr
		}
	}()

	// Close results when done
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
