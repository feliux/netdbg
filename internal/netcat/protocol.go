package netcat

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/feliux/netdbg/internal/logger"
)

// Connector is the interface that defines methods for connecting and listening.
type Connector interface {
	Connect(address string, port int, zero bool) error
	Listen(address string, port int) error
}

// TCPConnector implements the Connector interface for TCP.
type TCPConnector struct{}

// Connect establishes a TCP connection to the server.
func (c *TCPConnector) Connect(address string, port int, zero bool) error {
	hostPort := fmt.Sprintf("%s:%d", address, port)
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error("failed to close connection", "error", err)
		}
	}()

	if zero {
		logger.Info("zero mode invoked. Connection established.", "protocol", "tcp", "address", address, "port", port)
		return nil
	}

	logger.Info("connection established", "protocol", "tcp", "address", address, "port", port)
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		logger.Error("connection error", "protocol", "tcp", "address", address, "port", "error", err)
		return fmt.Errorf("connection error: %w", err)
	}

	logger.Info("connection closed successfully", "protocol", "tcp", "address", address, "port", port)
	return nil
}

// Listen starts a TCP server to listen for incoming connections.
func (c *TCPConnector) Listen(address string, port int) error {
	hostPort := fmt.Sprintf("%s:%d", address, port)
	logger.Info("starting TCP server", "address", address, "port", port)

	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		logger.Error("failed to start TCP server", "address", address, "port", port, "error", err)
		return err
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			logger.Error("failed to close listener", "error", err)
		}
	}()

	logger.Info("tcp server listening", "address", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("error accepting connection", "error", err)
			continue
		}
		logger.Info("accepted connection", "remote_address", conn.RemoteAddr().String())
		go processClient(conn)
	}
}

// UDPConnector implements the Connector interface for UDP.
type UDPConnector struct{}

// Connect establishes a UDP connection to the server.
func (c *UDPConnector) Connect(address string, port int, zero bool) error {
	hostPort := fmt.Sprintf("%s:%d", address, port)
	conn, err := net.Dial("udp", hostPort)
	if err != nil {
		logger.Error("failed to connect", "protocol", "udp", "address", address, "port", "error", err)
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error("failed to close connection", "error", err)
		}
	}()

	if zero {
		logger.Info("zero mode invoked. Connection established.", "protocol", "udp", "address", address, "port", port)
		return nil
	}

	logger.Info("connection established", "protocol", "udp", "address", address, "port", port)
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		logger.Error("connection error", "protocol", "udp", "address", address, "port", "error", err)
		return fmt.Errorf("connection error: %w", err)
	}

	logger.Info("connection closed successfully", "protocol", "udp", "address", address, "port", port)
	return nil
}

// Listen starts a UDP server to listen for incoming connections.
func (c *UDPConnector) Listen(address string, port int) error {
	// tbd implement udp logic
	return nil
}

// NewConnector creates a new instance of Connector based on the protocol.
func NewConnector(protocol string) Connector {
	switch protocol {
	case "udp":
		return &UDPConnector{}
	default:
		return &TCPConnector{}
	}
}

// processClient processes the data sent by a client.
func processClient(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error("failed to close client connection", "error", err)
		}
	}()

	logger.Info("processing client data", "remote_address", conn.RemoteAddr().String())
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		logger.Error("error processing client data", "remote_address", conn.RemoteAddr().String(), "error", err)
	}
	logger.Info("finished processing client data", "remote_address", conn.RemoteAddr().String())
}
