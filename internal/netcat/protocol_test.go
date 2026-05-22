package netcat

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestTCPConnector_Connect(t *testing.T) {
	setupLogger()
	// Start a mock TCP server
	address := "localhost"
	port := 5001
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
		if err != nil {
			t.Fatalf("failed to start mock server: %v", err)
		}
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("failed to accept connection: %v", err)
			return
		}
		defer conn.Close()

		// Simulate server response
		conn.Write([]byte("hello, client"))
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Test the TCPConnector's Connect method
	connector := &TCPConnector{}
	err := connector.Connect(address, port, false)
	if err != nil {
		t.Errorf("connect failed: %v", err)
	}
}

func TestTCPConnector_Listen(t *testing.T) {
	setupLogger()
	// Start the TCPConnector's Listen method in a goroutine
	address := "localhost"
	port := 5002
	connector := &TCPConnector{}
	go func() {
		err := connector.Listen(context.Background(), address, port)
		if err != nil {
			t.Errorf("listen failed: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Connect to the server as a client
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	message := "hello, server"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Errorf("failed to send data: %v", err)
	}
}

func TestUDPConnector_Connect(t *testing.T) {
	// Start a mock UDP server
	address := "localhost"
	port := 5003
	go func() {
		conn, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", address, port))
		if err != nil {
			t.Fatalf("failed to start mock UDP server: %v", err)
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		_, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			t.Errorf("failed to read from UDP client: %v", err)
			return
		}

		// Simulate server response
		conn.WriteTo([]byte("hello, udp client"), addr)
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Test the UDPConnector's Connect method
	connector := &UDPConnector{}
	err := connector.Connect(address, port, false)
	if err != nil {
		t.Errorf("connect failed: %v", err)
	}
	setupLogger()
}

func TestUDPConnector_Listen(t *testing.T) {
	setupLogger()
	// Start the UDPConnector's Listen method in a goroutine
	address := "localhost"
	port := 5004
	connector := &UDPConnector{}
	go func() {
		err := connector.Listen(context.Background(), address, port)
		if err != nil {
			t.Errorf("listen failed: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Connect to the server as a client
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", address, port))
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
}
