package netcat

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/feliux/netdbg/internal/logger"
)

// Connector is the interface that defines methods for connecting and listening.
type Connector interface {
	Connect(address string, port int, zero bool) error
	Listen(ctx context.Context, address string, port int) error
}

// TCPConnector implements the Connector interface for TCP.
type TCPConnector struct{}

// Connect establishes a TCP connection to the server.
func (c *TCPConnector) Connect(address string, port int, zero bool) error {
	hostPort := net.JoinHostPort(address, strconv.Itoa(port))
	logger.Info("connecting to remote endpoint", "protocol", "tcp", "address", address, "port", port)
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Debug("failed to close network connection", "error", err)
		}
	}()

	if zero {
		logger.Debug("zero-I/O mode: connection established", "protocol", "tcp", "address", address, "port", port)
		return nil
	}

	logger.Debug("connection established to remote endpoint", "protocol", "tcp", "address", address, "port", port)
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		logger.Debug("connection error while sending data", "protocol", "tcp", "address", address, "port", "error", err)
		return fmt.Errorf("connection error: %w", err)
	}

	logger.Debug("connection closed cleanly", "protocol", "tcp", "address", address, "port", port)
	return nil
}

// Listen starts a TCP server to listen for incoming connections, and supports context cancellation.
func (c *TCPConnector) Listen(ctx context.Context, address string, port int) error {
	hostPort := net.JoinHostPort(address, strconv.Itoa(port))
	logger.Info("starting TCP listener", "address", address, "port", port)

	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		logger.Debug("failed to start TCP listener", "address", address, "port", port, "error", err)
		return err
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			logger.Debug("failed to close TCP listener", "error", err)
		}
	}()

	logger.Info("TCP listener active", "address", listener.Addr().String())
	for {
		select {
		case <-ctx.Done():
			logger.Debug("listener context canceled, shutting down")
			return nil
		default:
			// Only set deadline if listener is a *net.TCPListener
			if tcpListener, ok := listener.(*net.TCPListener); ok {
				tcpListener.SetDeadline(time.Now().Add(200 * time.Millisecond))
			}
			conn, err := listener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue // check context again
				}
				logger.Debug("failed to accept connection", "error", err)
				continue
			}
			logger.Debug("accepted client connection", "remote_address", conn.RemoteAddr().String())
			go processClient(conn)
		}
	}
}

// UDPConnector implements the Connector interface for UDP.
type UDPConnector struct{}

// Connect establishes a UDP connection to the server.
func (c *UDPConnector) Connect(address string, port int, zero bool) error {
	hostPort := net.JoinHostPort(address, strconv.Itoa(port))
	logger.Info("connecting to remote endpoint", "protocol", "udp", "address", address, "port", port)
	conn, err := net.Dial("udp", hostPort)
	if err != nil {
		logger.Debug("failed to connect to remote endpoint", "protocol", "udp", "address", address, "port", "error", err)
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Debug("failed to close network connection", "error", err)
		}
	}()

	if zero {
		logger.Debug("zero-I/O mode: connection established", "protocol", "udp", "address", address, "port", port)
		return nil
	}

	logger.Debug("connection established to remote endpoint", "protocol", "udp", "address", address, "port", port)
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		logger.Debug("connection error while sending data", "protocol", "udp", "address", address, "port", "error", err)
		return fmt.Errorf("connection error: %w", err)
	}

	logger.Debug("connection closed cleanly", "protocol", "udp", "address", address, "port", port)
	return nil
}

// Listen starts a UDP server to listen for incoming connections.
func (c *UDPConnector) Listen(ctx context.Context, address string, port int) error {
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
			logger.Debug("failed to close client connection cleanly", "error", err)
		}
	}()

	logger.Debug("processing client data stream", "remote_address", conn.RemoteAddr().String())
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		logger.Debug("error processing client data stream", "remote_address", conn.RemoteAddr().String(), "error", err)
	}
	logger.Debug("finished processing client data stream", "remote_address", conn.RemoteAddr().String())
}
