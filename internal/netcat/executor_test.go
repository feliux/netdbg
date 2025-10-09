package netcat

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

// getFreePort find a free TCP port for testing listen/connect
func getFreePort(t *testing.T) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func TestDefaultExecutor_Connect_Success(t *testing.T) {
	setupLogger()
	port := getFreePort(t)

	// Start a dummy TCP server
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("failed to start dummy server: %v", err)
	}
	defer ln.Close()

	done := make(chan struct{})
	go func() {
		conn, err := ln.Accept()
		if err == nil {
			conn.Close()
		}
		close(done)
	}()

	opts := &Options{
		Address:  "127.0.0.1",
		Port:     port,
		Protocol: "tcp",
		Listen:   false,
		Zero:     true,
		Timeout:  2 * time.Second,
	}
	executor := &DefaultExecutor{}
	result := executor.Execute(context.Background(), opts)
	if result.Error != nil {
		t.Errorf("expected no error, got: %v", result.Error)
	}
	if !result.Success {
		t.Errorf("expected success, got: %v", result.Success)
	}
	<-done
}

func TestDefaultExecutor_Connect_Fail(t *testing.T) {
	setupLogger()
	// Use a port unlikely to be open
	opts := &Options{
		Address:  "127.0.0.1",
		Port:     65000,
		Protocol: "tcp",
		Listen:   false,
		Zero:     true,
		Timeout:  1 * time.Second,
	}
	executor := &DefaultExecutor{}
	result := executor.Execute(context.Background(), opts)
	if result.Error == nil {
		t.Errorf("expected error, got nil")
	}
	if result.Success {
		t.Errorf("expected failure, got success")
	}
}

func TestDefaultExecutor_ListenAndConnect(t *testing.T) {
	setupLogger()
	port := getFreePort(t)
	listenOpts := &Options{
		Address:  "127.0.0.1",
		Port:     port,
		Protocol: "tcp",
		Listen:   true,
		Zero:     false,
		Timeout:  2 * time.Second,
	}
	connectOpts := &Options{
		Address:  "127.0.0.1",
		Port:     port,
		Protocol: "tcp",
		Listen:   false,
		Zero:     true,
		Timeout:  2 * time.Second,
	}
	executor := &DefaultExecutor{}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		// Redirect os.Stdout temporarily to avoid polluting test output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		result := executor.Execute(ctx, listenOpts)

		w.Close()
		os.Stdout = oldStdout
		// Drain the pipe
		_, _ = io.ReadAll(r)
		r.Close()

		if result.Error != nil {
			t.Errorf("listen error: %v", result.Error)
		}
		close(done)
	}()

	// Give the listener a moment to start
	time.Sleep(100 * time.Millisecond)

	result := executor.Execute(context.Background(), connectOpts)
	if result.Error != nil {
		t.Errorf("connect error: %v", result.Error)
	}

	// Cancel the listener context to shut down the server
	cancel()
	<-done
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
