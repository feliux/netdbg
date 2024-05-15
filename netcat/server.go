package netcat

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Start a TCP server for listening connections
func StartServer(addr string, port int) {
	hostPort := fmt.Sprintf("%s:%d", addr, port)
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		panic(err)
	}
	log.Printf("listening for connections on %s", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection from client: %s", err)
		} else {
			go processClient(conn)
		}
	}
}

// Process data sent by client
func processClient(conn net.Conn) {
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
}
