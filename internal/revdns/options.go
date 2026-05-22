package revdns

// Options holds the configuration for a reverse DNS lookup operation.
type Options struct {
	Addr       string // IP address to perform reverse lookup on (single IP)
	ResolverIP string // Custom DNS resolver IP (optional)
	Protocol   string // Protocol to use with custom resolver ("udp" or "tcp")
	Port       int    // Port for the DNS resolver
	Threads    int    // Number of concurrent workers
	DomainOnly bool   // Output only domain names (no IP)
	File       string // Path to file with list of IPs (optional)
}
