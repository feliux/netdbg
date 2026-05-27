package revdns

import (
	"context"
	"sync"

	"github.com/feliux/netdbg/internal/logger"
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

	source := "address"
	if opts.File != "" {
		source = "file"
	}
	resolver := opts.ResolverIP
	if resolver == "" {
		resolver = "system"
	}
	logger.Info("starting reverse DNS resolution", "source", source, "threads", opts.Threads, "resolver", resolver, "protocol", opts.Protocol, "port", opts.Port, "domain_only", opts.DomainOnly)

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
			logger.Info("loading IP addresses from file", "file", opts.File)
			if err := feedFromFile(opts.File, work); err != nil {
				results <- Result{Error: err}
			}
		} else if opts.Addr != "" {
			logger.Info("looking up IP address", "ip", opts.Addr)
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
