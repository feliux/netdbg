package netcat

import (
	"context"
	"time"

	"github.com/spf13/cobra"
)

// ExecuteCommand builds options from flags and runs the netcat executor.
func ExecuteCommand(cmd *cobra.Command, args []string) Result {
	listen, _ := cmd.Flags().GetBool("listen")
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")
	protocol, _ := cmd.Flags().GetString("protocol")
	zero, _ := cmd.Flags().GetBool("zero")
	// Optional: add timeout flag in the future
	timeout := 10 * time.Second

	opts := &Options{
		Address:  address,
		Port:     port,
		Protocol: protocol,
		Listen:   listen,
		Zero:     zero,
		Timeout:  timeout,
	}

	executor := &DefaultExecutor{}
	result := executor.Execute(context.Background(), opts)

	return result
}
