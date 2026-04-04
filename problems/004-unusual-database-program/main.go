package main

import (
	"log"
	"net"
	"strings"

	"github.com/hellinvoid/network-programming/problems/004-unusual-database-program/service"
)

const (
	PROTOCOL        = "udp"
	ADDRESS         = ":6969"
	MAX_BUFFER_SIZE = 1000
)

func main() {

	// Convert string to net.UDPAddr
	addr, err := net.ResolveUDPAddr(PROTOCOL, ADDRESS)
	if err != nil {
		panic(err)
	}

	// Listen for UDP packets
	udp, err := net.ListenUDP(PROTOCOL, addr)
	if err != nil {
		panic(err)
	}
	defer udp.Close()

	// Instantiate a new db service
	db := service.NewDatabase()

	log.Printf("Listening for UDP on %s\n", addr.String())
	// Handle all udp messages in sequence NO CONCURRENCY
	for {
		handleUDP(udp, db)
	}

}

func handleUDP(udp *net.UDPConn, db service.Database) {
	buf := make([]byte, MAX_BUFFER_SIZE)
	// Read the message
	n, addr, err := udp.ReadFromUDP(buf)
	if err != nil {
		return
	}
	// Forward it to further checks
	handlePacket(buf[:n], addr, udp, db)
}
func handlePacket(buf []byte, addr *net.UDPAddr, udp *net.UDPConn, db service.Database) {
	str := string(buf)

	// Split the message for key value pair
	before, after, ok := strings.Cut(str, "=")

	// Insert if key value pair and retrieve if not
	if ok {
		key := before
		value := after
		db.Insert(key, value)
	} else {
		value := db.Retrieve(str)
		udp.WriteToUDP([]byte(value), addr)
	}

}
