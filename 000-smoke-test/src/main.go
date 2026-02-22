package main

import (
	"io"
	"net"
)

const (
	PROTOCOL    = "tcp"
	ADDRESS     = ":6969"
	// BUFFER_SIZE = 4096
)

func main() {
	listner, err := net.Listen(PROTOCOL, ADDRESS)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listner.Accept()
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
