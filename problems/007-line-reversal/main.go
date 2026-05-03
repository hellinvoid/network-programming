package main

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/hellinvoid/network-programming/problems/007-line-reversal/dto"
	"github.com/hellinvoid/network-programming/problems/007-line-reversal/service"
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

	log.Printf("Listening for UDP on %s\n", addr.String())
	// Handle all udp messages in sequence NO CONCURRENCY

	lrcp := service.NewLRCP(udp)
	handleLRCP(lrcp)
}

func handleLRCP(lrcp *service.LRCP) {

	// Message routing happens here
	buf := make([]byte, MAX_BUFFER_SIZE)
	for {
		n, addr, err := lrcp.Udp.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		str := string(buf[:n])
		log.Printf("From : %s\n Msg: %s", addr.String(), str)
		if str[0] != '/' && str[n-1] != '/' {
			continue
		}

		msgType, aft, found := strings.Cut(str[1:n], "/")
		if !found {
			continue
		}

		sessionIdStr, aft, found := strings.Cut(aft, "/")
		if !found {
			continue
		}

		sessionId, err := strconv.Atoi(sessionIdStr)
		if err != nil {
			continue
		}

		msg := &dto.Message{
			Addr:      addr,
			MsgType:   msgType,
			SessionId: uint64(sessionId),
			Buf:       []byte(aft),
		}

		lrcp.SendToSession(msg)
	}

}
