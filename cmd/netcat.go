package cmd

import (
	"log/slog"

	"github.com/feliux/netdbg/internal/netcat"
	"github.com/spf13/cobra"
)

var netcatCmd = &cobra.Command{
	Use:   "nc",
	Short: "Minimal netcat tool",
	Long:  "Connect to somewhere or listen for inbound connections.",
	Run: func(cmd *cobra.Command, args []string) {
		netcat.ExecuteCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(netcatCmd)
	netcatCmd.Flags().BoolP("listen", "l", false, "listen mode, for inbound connects")
	netcatCmd.Flags().StringP("address", "a", "localhost", "hostname to connect (client)")
	netcatCmd.Flags().IntP("port", "p", 5000, "port to connect (client) or bind (server)")
	netcatCmd.Flags().StringP("protocol", "P", "tcp", "protocol to use (tcp|udp)")
	netcatCmd.Flags().BoolP("zero", "z", false, "zero-I/O mode (used for scanning)")
	err := netcatCmd.MarkFlagRequired("address")
	if err != nil {
		slog.Error("can not require flag address", "err", err)
	}
	err = netcatCmd.MarkFlagRequired("port")
	if err != nil {
		slog.Error("can not require flag port", "err", err)
	}
}
