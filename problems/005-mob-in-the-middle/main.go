package main

import (
	"log"
	"net"

	"github.com/hellinvoid/network-programming/problems/005-mob-in-the-middle/service"
)

const (
	PROTOCOL                = "tcp"
	PROXY_ADDRESS           = ":6969"
	UPSTREAM_SERVER_ADDRESS = "chat.protohackers.com:16963"
)

func main() {
	listener, err := net.Listen(PROTOCOL, PROXY_ADDRESS)
	if err != nil {
		panic(err)
	}

	log.Printf("Proxy server listening on %s", listener.Addr().String())
	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(client net.Conn) {
	defer client.Close()

	// Open connection for upstream server
	server, err := net.Dial(PROTOCOL, UPSTREAM_SERVER_ADDRESS)
	if err != nil {
		return
	}
	defer server.Close()

	// Proxy the request from client to server
	go service.ProxyRequest(client, server)
	// Proxy the request from server to client
	service.ProxyRequest(server, client)
}
