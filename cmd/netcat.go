package cmd

import (
	"fmt"
	"os"

	"github.com/feliux/netdbg/internal/logger"
	"github.com/feliux/netdbg/internal/netcat"
	"github.com/spf13/cobra"
)

var netcatCmd = &cobra.Command{
	Use:   "nc",
	Short: "Minimal netcat tool",
	Long: `Minimal netcat tool.

Usage examples:
  netdbg nc -a <address> -p <port>
  netdbg nc --listen -a <address> -p <port>
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		if address == "" {
			err := cmd.Help()
			if err != nil {
				logger.Error("failed to execute help command", "err", err)
			}
			fmt.Fprintln(os.Stderr, "Error: you must specify the destination address using the -a or --address flag.")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		result := netcat.ExecuteCommand(cmd, args)
		if result.Error != nil {
			message := "unable to connect"
			hint := "check host, port, or connectivity"
			if result.ListenMode {
				message = "unable to start server"
				hint = "use another port or check permissions"
			}
			WriteError(os.Stderr, ErrorOutput{
				Command: "nc",
				Message: message,
				Cause:   result.Error,
				Hint:    hint,
			})
			return
		}
		fmt.Fprintln(os.Stdout, result.Message)
	},
}

func init() {
	rootCmd.AddCommand(netcatCmd)
	netcatCmd.Flags().BoolP("listen", "l", false, "listen mode, for inbound connects")
	netcatCmd.Flags().StringP("address", "a", "", "hostname to connect (client)")
	netcatCmd.Flags().IntP("port", "p", 5000, "port to connect (client) or bind (server)")
	netcatCmd.Flags().StringP("protocol", "P", "tcp", "protocol to use (tcp|udp)")
	netcatCmd.Flags().BoolP("zero", "z", false, "zero-I/O mode (used for scanning)")
}
