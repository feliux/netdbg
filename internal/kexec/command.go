package kexec

import (
	"context"

	"github.com/spf13/cobra"
)

// ExecuteCommand builds options from flags and runs the kexec executor.
func ExecuteCommand(cmd *cobra.Command, args []string) (*Options, error) {
	opts := ParseOptionsFromFlags(cmd.Flags())

	executor := &DefaultExecutor{}
	err := executor.Execute(context.Background(), opts, args)
	return opts, err
}
