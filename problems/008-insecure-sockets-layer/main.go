package main

import (
	"bufio"
	"bytes"
	"log"
	"net"

	"github.com/hellinvoid/network-programming/problems/008-insecure-sockets-layer/service"
)

const (
	PROTOCOL      = "tcp"
	PROXY_ADDRESS = ":6969"
)

func main() {
	listener, err := net.Listen(PROTOCOL, PROXY_ADDRESS)
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

	r := bufio.NewReader(conn)

	cipher, err := r.ReadBytes(service.END)
	if err != nil || len(cipher) < 2 {
		return
	}
	secLast := cipher[len(cipher)-2]

	// 0 is a valid N 
	if secLast == service.ADD_N || secLast == service.XOR_N {
		r.ReadByte()
	} else {
		cipher = cipher[:len(cipher)-1]
	}


	rec := 0
	sent := 0
	for {
		decoded := make([]byte, 0)
		original := make([]byte, 0)
		for {
			b, err := r.ReadByte()
			if err != nil {
				return
			}
			d := service.Decode(b, cipher, 0, rec)
			rec++
			if d == '\n' {
				break
			}
			original = append(original, b)
			decoded = append(decoded, d)
		}
		if bytes.Equal(original, decoded) {
			return
		}


		toy := service.GetToy(string(decoded)) + "\n"

		encoded := make([]byte, len(toy))
		for i, d := range []byte(toy) {
			encoded[i] = service.Encode(d, cipher, 0, sent)
			sent++
		}

		conn.Write(encoded)
	}

}
