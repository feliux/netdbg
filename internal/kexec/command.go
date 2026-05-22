package kexec

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func ExecuteCommand(cmd *cobra.Command, args []string) {
	opts := ParseOptionsFromFlags(cmd.Flags())

	executor := &DefaultExecutor{}
	if err := executor.Execute(context.Background(), opts, args); err != nil {
		fmt.Fprintln(os.Stderr, "kexec error:", err)
		os.Exit(1)
	}
}
