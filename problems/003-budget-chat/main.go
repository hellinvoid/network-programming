package main

import (
	"log"
	"net"

	"github.com/hellinvoid/network-programming/problems/003-budget-chat/service"
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

	cr := service.NewChatRoom()

	log.Printf("Server listening on %s", listener.Addr().String())
	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn, cr)
	}

}

func handleConn(conn net.Conn, cr service.ChatRoom) {
	defer conn.Close()

	cl, err := service.NewClient(conn, cr)
	if err != nil {
		return
	}
	
	cl.Start()
}
