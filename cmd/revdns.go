// Package cmd revdns.go is a Cobra cli entrypoint for reverse DNS lookups.
package cmd

import (
	"fmt"
	"os"

	"github.com/feliux/netdbg/internal/logger"
	"github.com/feliux/netdbg/internal/revdns"

	"github.com/spf13/cobra"
)

var revdnsCmd = &cobra.Command{
	Use:   "revdns",
	Short: "Reverse DNS lookup",
	Long: `Reverse DNS lookup tool.

You can perform a reverse DNS lookup for a single IP address or a list of IPs from a file.

Usage examples:
  netdbg revdns -a <ip> -p <port> -r <resolver_ip> -P <udp|tcp>
  netdbg revdns -f <file_with_ips> -p <port> -r <resolver_ip> -P <udp|tcp>

You must specify the address using the -a or --address flag, or provide a file with the -f or --file flag.
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("address")
		file, _ := cmd.Flags().GetString("file")
		if addr == "" && file == "" {
			err := cmd.Help()
			if err != nil {
				logger.Error("can not execute help command", "err", err)
			}
			fmt.Fprintln(os.Stderr, "Error: you must specify either the address with -a/--address or a file with -f/--file.")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		revdns.ExecuteCommand(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(revdnsCmd)
	revdnsCmd.Flags().StringP("address", "a", "", "address for reverse DNS lookups")
	revdnsCmd.Flags().StringP("resolver", "r", "", "IP of the DNS resolver to use for lookups")
	revdnsCmd.Flags().StringP("protocol", "P", "udp", "protocol to use for lookups (udp | tcp)")
	revdnsCmd.Flags().StringP("file", "f", "", "file containing a list of IPs for lookup")
	revdnsCmd.Flags().IntP("port", "p", 53, "port to bother the specified DNS resolver on")
	revdnsCmd.Flags().IntP("threads", "t", 5, "threads to use when reversing from file")
	revdnsCmd.Flags().BoolP("domain", "d", false, "output only domain names")

}
