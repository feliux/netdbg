package netcat

import (
	"time"
)

// Options holds the configuration for a netcat operation.
type Options struct {
	Address  string        // Hostname or IP to connect or bind to
	Port     int           // Port to connect or bind
	Protocol string        // Protocol to use ("tcp" or "udp")
	Listen   bool          // Listen mode (server) if true, connect mode (client) if false
	Zero     bool          // Zero-I/O mode (used for scanning)
	Timeout  time.Duration // Optional timeout for the operation
}
