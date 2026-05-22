package revdns

// Result represents the outcome of a reverse DNS lookup.
type Result struct {
	IP     string // The IP address that was looked up
	Domain string // The resolved domain name (if any)
	Error  error  // Any error encountered during the lookup
}
