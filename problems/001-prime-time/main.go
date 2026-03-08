package main

import (
	"encoding/json"
	"net"
	"time"

	"github.com/hellinvoid/network-programming/problems/001-prime-time/dto"
	"github.com/hellinvoid/network-programming/problems/001-prime-time/service"
)

const (
	PROTOCOL = "tcp"
	ADDRESS  = ":6969"
)

func main() {
	// Create a listner
	listener, err := net.Listen(PROTOCOL, ADDRESS)
	if err != nil {
		panic(err)
	}

	// Accept connections
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

	// Initialize encoder and decoder
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	for {
		// Setting a read deadline to handle half correct json cases [IMPORTANT]
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		var req dto.Request
		var res dto.Response

		err := dec.Decode(&req)

		// Return empty if malformed request
		if err != nil || req.Number == nil || req.Method != "isPrime" {
			res.Method = "BOO"
			enc.Encode(res)
			return
		}

		res = service.IsPrime(&req)
		enc.Encode(res)
	}
}
