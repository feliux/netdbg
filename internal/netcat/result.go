package netcat

// Result represents the outcome of a netcat operation.
type Result struct {
	Success       bool   // Indicates if the operation was successful
	Error         error  // Error encountered during the operation, if any
	Message       string // Optional message with additional info
	ListenMode    bool   // True if operation was in listen mode
	RemoteAddr    string // Remote address connected to or accepted from
	LocalAddr     string // Local address used for the operation
	BytesSent     int    // Number of bytes sent (if applicable)
	BytesReceived int    // Number of bytes received (if applicable)
}
