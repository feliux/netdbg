// Package cmd root.go is a Cobra cli entrypoint.
package cmd

import (
	"os"

	"github.com/feliux/netdbg/internal/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "netdbg",
	Short: "Net debugger CLI",
	Long:  `Set of tools for testing and debugging connectivity issues.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		verbose, _ := cmd.Flags().GetBool("verbose")
		levelFlag, _ := cmd.Flags().GetString("log-level")
		formatFlag, _ := cmd.Flags().GetString("log-format")

		if verbose {
			levelFlag = "info"
			formatFlag = logger.FormatText
		}
		if debug {
			levelFlag = "debug"
			formatFlag = logger.FormatJSON
		}

		level, err := logger.ParseLevel(levelFlag)
		if err != nil {
			return err
		}
		normalizedFormat, err := logger.NormalizeFormat(formatFlag)
		if err != nil {
			return err
		}
		logger.InitLogger(level, os.Stderr, normalizedFormat)
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		WriteError(os.Stderr, ErrorOutput{
			Message: "command execution failed",
			Cause:   err,
			Hint:    "use --help for usage information",
		})
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose logging (text info)")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging (overrides --log-level)")
	rootCmd.PersistentFlags().String("log-level", "warn", "set log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-format", logger.FormatText, "log format (json|text)")
}
