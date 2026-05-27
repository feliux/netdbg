package revdns

import (
	"context"

	"github.com/spf13/cobra"
)

// ExecuteCommand builds options from flags and starts reverse DNS execution.
func ExecuteCommand(cmd *cobra.Command, args []string) (*Options, <-chan Result) {
	addr, _ := cmd.Flags().GetString("address")
	resolverIP, _ := cmd.Flags().GetString("resolver")
	protocol, _ := cmd.Flags().GetString("protocol")
	file, _ := cmd.Flags().GetString("file")
	port, _ := cmd.Flags().GetInt("port")
	threads, _ := cmd.Flags().GetInt("threads")
	domain, _ := cmd.Flags().GetBool("domain")

	opts := &Options{
		Addr:       addr,
		ResolverIP: resolverIP,
		Protocol:   protocol,
		Port:       port,
		Threads:    threads,
		DomainOnly: domain,
		File:       file,
	}

	executor := &DefaultExecutor{}
	ctx := context.Background()
	results := executor.Execute(ctx, opts)

	return opts, results
}
