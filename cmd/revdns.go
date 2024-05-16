package cmd

import (
	"os"

	"github.com/feliux/netdbg/revdns"
	"github.com/spf13/cobra"
)

// revdnsCmd represents the reverse lookup command
var revdnsCmd = &cobra.Command{
	Use:   "revdns",
	Short: "Reverse DNS lookup",
	Long: `
Make a reverse DNS lookup:  netdbg revdns -a [ip] -p [port] -r [resolver_ip] -P [udp|tcp] -d
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("address")
		file, _ := cmd.Flags().GetString("file")
		if addr == "" && file == "" {
			cmd.Help()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("address")
		resolverIP, _ := cmd.Flags().GetString("resolver")
		protocol, _ := cmd.Flags().GetString("protocol")
		file, _ := cmd.Flags().GetString("file")
		port, _ := cmd.Flags().GetInt("port")
		threads, _ := cmd.Flags().GetInt("threads")
		domain, _ := cmd.Flags().GetBool("domain")
		revdns.RevDNS(addr, resolverIP, protocol, file, port, threads, domain)
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
	// revdnsCmd.MarkFlagRequired("address")

}
