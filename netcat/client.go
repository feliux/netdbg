package netcat

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Start a TCP client. Include zero mode: https://unix.stackexchange.com/questions/589561/what-is-nc-z-used-for
func StartClient(addr string, port int, zero bool) {
	hostPort := fmt.Sprintf("%s:%d", addr, port)
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		log.Printf("can not connect to server: %s\n", err)
		return
	} else if err == nil && zero {
		fmt.Println("zero mode invoked. Connection established.")
		conn.Close()
		return
	}
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		fmt.Printf("connection error: %s\n", err)
	}
}
