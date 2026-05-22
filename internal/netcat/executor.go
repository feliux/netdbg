package netcat

import (
	"context"
	"fmt"
)

// Executor defines the interface for executing netcat operations.
type Executor interface {
	Execute(ctx context.Context, opts *Options) Result
}

// DefaultExecutor implements Executor with standard netcat logic.
type DefaultExecutor struct{}

// Execute performs the netcat operation (connect or listen) based on the provided options.
func (e *DefaultExecutor) Execute(ctx context.Context, opts *Options) Result {
	connector := NewConnector(opts.Protocol)
	if opts.Listen {
		err := connector.Listen(ctx, opts.Address, opts.Port)
		if err != nil {
			return Result{Success: false, Error: fmt.Errorf("listen error: %w", err), ListenMode: true}
		}
		return Result{Success: true, ListenMode: true, Message: "listening..."}
	} else {
		err := connector.Connect(opts.Address, opts.Port, opts.Zero)
		if err != nil {
			return Result{Success: false, Error: fmt.Errorf("connect error: %w", err)}
		}
		if opts.Zero {
			return Result{Success: true, Message: "zero-I/O mode: connection established and closed"}
		}
		return Result{Success: true, Message: "connection established"}
	}
}
