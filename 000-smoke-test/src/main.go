package main

import (
	"io"
	"net"
)

const (
	PROTOCOL = "tcp"
	ADDRESS  = ":6969"
	// BUFFER_SIZE = 4096
)

func main() {
	listener, err := net.Listen(PROTOCOL, ADDRESS)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	io.Copy(conn, conn)
}
