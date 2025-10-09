// Package revdns for reverse DNS lookups.
package revdns

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
)

type RevDNSLookup struct {
	Addr, ResolverIP, Protocol string
	Port, Threads              int
	Domain                     bool
}

var (
	opts *RevDNSLookup
)

// RevDNS build the worker pool for a concurrent reverse lookup.
func RevDNS(addr, resolverIP, protocol, file string, port, threads int, domain bool) {
	var wg sync.WaitGroup
	opts = &RevDNSLookup{
		Addr:       addr,
		ResolverIP: resolverIP,
		Protocol:   protocol,
		Threads:    threads,
		Port:       port,
		Domain:     domain,
	}
	work := make(chan string)
	if file != "" {
		go readFile(&work, file)
	} else {
		go send(&work, opts.Addr)
	}
	// go func(filePath string) {
	// 	file, err := os.Open(filePath)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer file.Close()
	// 	sc := bufio.NewScanner(file)
	// 	for sc.Scan() {
	// 		work <- sc.Text() // GET the line string
	// 	}
	// 	close(work)
	// }(file)
	for i := 0; i < opts.Threads; i++ {
		wg.Add(1)
		go doRev(work, &wg, opts)
	}
	wg.Wait()
}

// send sends an ip to a worker.
func send(work *chan string, ip string) {
	*work <- ip
	close(*work)
}

// readFile reads a file withs IPs and send it to a worker 1 by 1.
func readFile(work *chan string, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	// defer file.Close()
	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("can not close file", "err", err)
		}
	}()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		*work <- sc.Text() // GET the line string
	}
	close(*work)
}

// doRev does the reverse DNS lookup.
func doRev(work chan string, wg *sync.WaitGroup, opts *RevDNSLookup) {
	defer wg.Done()
	var r *net.Resolver
	if opts.ResolverIP != "" {
		r = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, opts.Protocol, fmt.Sprintf("%s:%d", opts.ResolverIP, opts.Port))
			},
		}
	}
	for ip := range work {
		addrs, err := r.LookupAddr(context.Background(), ip)
		if err != nil {
			fmt.Println("error occurred:", err)
			continue
		}
		for _, a := range addrs {
			if opts.Domain {
				fmt.Println(strings.TrimRight(a, "."))
			} else {
				fmt.Println(ip, "\t", a)
			}
		}
	}
}
