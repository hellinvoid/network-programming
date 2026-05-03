package dto

import "net"

type Message struct {
	Addr      *net.UDPAddr
	MsgType   string
	SessionId uint64
	Buf       []byte
}
