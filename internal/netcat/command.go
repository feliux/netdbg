package netcat

import (
	"context"
	"fmt"
	"time"

	"github.com/feliux/netdbg/internal/logger"
	"github.com/spf13/cobra"
)

func ExecuteCommand(cmd *cobra.Command, args []string) {
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

	if result.Error != nil {
		if listen {
			logger.Error("error starting server", "err", result.Error)
		} else {
			logger.Error("error connecting to server", "err", result.Error)
		}
		fmt.Println(result.Error)
	} else {
		fmt.Println(result.Message)
	}
}
