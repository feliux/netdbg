package revdns

import (
	"context"
	"fmt"
	"os"

	"github.com/feliux/netdbg/internal/logger"
	"github.com/spf13/cobra"
)

func ExecuteCommand(cmd *cobra.Command, args []string) {
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
	for result := range executor.Execute(ctx, opts) {
		if result.Error != nil {
			logger.Error("reverse DNS lookup error", "ip", result.IP, "err", result.Error)
			fmt.Fprintf(os.Stderr, "error: %v\n", result.Error)
			continue
		}
		if opts.DomainOnly {
			fmt.Println(result.Domain)
		} else {
			fmt.Printf("%s\t%s\n", result.IP, result.Domain)
		}
	}
}
