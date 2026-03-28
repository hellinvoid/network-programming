package main

import (
	"log"
	"net"

	"github.com/hellinvoid/network-programming/problems/002-means-to-an-end/dto"
	"github.com/hellinvoid/network-programming/problems/002-means-to-an-end/service"
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
	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// Create separate asset for each connection
	asset := service.NewAsset()
	for {

		// Create a valid request
		req, err := dto.NewRequest(conn)
		if err != nil {
			return
		}

		res, err := asset.HandleRequest(req)
		if err != nil {
			return
		}
		if res == nil {
			continue
		}
		
		err = service.SendResponse(*res, conn)
		if err != nil {
			return
		}
	}

}
