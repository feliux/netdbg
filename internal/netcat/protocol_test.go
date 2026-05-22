package netcat

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestTCPConnector_Connect(t *testing.T) {
	setupLogger()
	// Start a mock TCP server
	address := "localhost"
	port := 5001
	errCh := make(chan error, 1)

	go func() {
		listener, err := net.Listen("tcp", net.JoinHostPort(address, strconv.Itoa(port)))
		if err != nil {
			errCh <- fmt.Errorf("failed to start mock server: %w", err)
			return
		}
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			errCh <- fmt.Errorf("failed to accept connection: %w", err)
			return
		}
		defer conn.Close()

		// Simulate server response
		if _, err := conn.Write([]byte("hello, client")); err != nil {
			errCh <- fmt.Errorf("failed to write response: %w", err)
			return
		}
		errCh <- nil
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Test the TCPConnector's Connect method
	connector := &TCPConnector{}
	err := connector.Connect(address, port, false)
	if err != nil {
		t.Errorf("connect failed: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("server error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Errorf("server did not respond in time")
	}
}

func TestTCPConnector_Listen(t *testing.T) {
	setupLogger()
	// Start the TCPConnector's Listen method in a goroutine
	address := "localhost"
	port := 5002
	connector := &TCPConnector{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- connector.Listen(ctx, address, port)
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Connect to the server as a client
	conn, err := net.Dial("tcp", net.JoinHostPort(address, strconv.Itoa(port)))
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	message := "hello, server"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Errorf("failed to send data: %v", err)
	}

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("listen failed: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Errorf("listener did not stop in time")
	}
}

func TestUDPConnector_Connect(t *testing.T) {
	// Start a mock UDP server
	address := "127.0.0.1"
	port := 5003
	errCh := make(chan error, 1)

	originalStdin := os.Stdin
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}
	os.Stdin = reader
	defer func() {
		os.Stdin = originalStdin
		reader.Close()
	}()

	go func() {
		conn, err := net.ListenPacket("udp", net.JoinHostPort(address, strconv.Itoa(port)))
		if err != nil {
			errCh <- fmt.Errorf("failed to start mock UDP server: %w", err)
			return
		}
		defer conn.Close()

		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buffer := make([]byte, 1024)
		_, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			errCh <- fmt.Errorf("failed to read from UDP client: %w", err)
			return
		}

		// Simulate server response
		if _, err := conn.WriteTo([]byte("hello, udp client"), addr); err != nil {
			errCh <- fmt.Errorf("failed to write response: %w", err)
			return
		}
		errCh <- nil
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	go func() {
		_, _ = writer.Write([]byte("hello, udp server"))
		writer.Close()
	}()

	// Test the UDPConnector's Connect method
	connector := &UDPConnector{}
	err = connector.Connect(address, port, false)
	if err != nil {
		t.Errorf("connect failed: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("server error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Errorf("server did not respond in time")
	}
	setupLogger()
}

func TestUDPConnector_Listen(t *testing.T) {
	setupLogger()
	// Start the UDPConnector's Listen method in a goroutine
	address := "127.0.0.1"
	port := 5004
	connector := &UDPConnector{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- connector.Listen(ctx, address, port)
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Connect to the server as a client
	conn, err := net.Dial("udp", net.JoinHostPort(address, strconv.Itoa(port)))
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Simulate client sending data
	message := "Hello, UDP server!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Errorf("failed to send data: %v", err)
	}

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("listen failed: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Errorf("listener did not stop in time")
	}
}
