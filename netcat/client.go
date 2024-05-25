// Package netcat client.go builds a TCP clients.
package netcat

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

// StartClient starts a TCP client. Include zero mode:
// https://unix.stackexchange.com/questions/589561/what-is-nc-z-used-for
func StartClient(addr string, port int, zero bool) {
	hostPort := fmt.Sprintf("%s:%d", addr, port)
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		slog.Error("can not connect to server", "err", err)
		return
	} else if err == nil && zero {
		fmt.Println("zero mode invoked. Connection established.")
		err = conn.Close()
		if err != nil {
			slog.Error("can not close connection", "err", err)
		}
		return
	}
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		fmt.Printf("connection error: %s\n", err)
	}
}
