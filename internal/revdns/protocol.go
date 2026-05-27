package revdns

import (
	"context"
	"fmt"
	"net"
	"strings"
)

// worker processes IPs from the work channel and emits results.
func (e *DefaultExecutor) worker(ctx context.Context, opts *Options, work <-chan string, results chan<- Result) {
	resolver := getResolver(opts)
	for ip := range work {
		addrs, err := resolver.LookupAddr(ctx, ip)
		if err != nil {
			results <- Result{IP: ip, Error: err}
			continue
		}
		for _, a := range addrs {
			domain := strings.TrimRight(a, ".")
			if opts.DomainOnly {
				results <- Result{Domain: domain}
			} else {
				results <- Result{IP: ip, Domain: domain}
			}
		}
	}
}

// getResolver returns a resolver based on options or the system default.
func getResolver(opts *Options) *net.Resolver {
	if opts.ResolverIP != "" {
		return &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, opts.Protocol, fmt.Sprintf("%s:%d", opts.ResolverIP, opts.Port))
			},
		}
	}
	return net.DefaultResolver
}
