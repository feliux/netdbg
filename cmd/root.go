// Package cmd root.go is a Cobra cli entrypoint.
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "netdbg",
	Short: "Net debugger CLI",
	Long:  `Set of tools for testing and debugging connectivity issues.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("error occurred calling root cmd", "err", err)
		os.Exit(1)
	}
}

func init() {}
