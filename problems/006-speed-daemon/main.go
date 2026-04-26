package main

import (
	"log"
	"net"

	"github.com/hellinvoid/network-programming/problems/006-speed-daemon/service"
)

const (
	PROTOCOL = "tcp"
	ADDRESS  = ":6969"
)

func main() {
	listener, err := net.Listen(PROTOCOL, ADDRESS)
	if err != nil {
		panic(err)
	}
	log.Printf("Server listening on %s", listener.Addr().String())

	// Creating a new ticket system to handle new requests
	ts := service.NewTicketSystem()
	go ts.HandlePlate()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// Let the ticket system handle to connections
		go ts.HandleConn(conn)
	}

}
