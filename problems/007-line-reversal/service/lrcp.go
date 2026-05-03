package service

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hellinvoid/network-programming/common"
	"github.com/hellinvoid/network-programming/problems/007-line-reversal/dto"
)

const (
	SESSION_EXPIRY_TIMEOUT = 60 * time.Second
	RETRANSMISSION_TIMEOUT = 3 * time.Second
	CHANNEL_BUFFER         = 500
)

type LRCP struct {
	Udp     *net.UDPConn
	Session map[uint64]*common.CloseWrapper[*dto.Message]
}

func NewLRCP(udp *net.UDPConn) *LRCP {
	return &LRCP{
		Session: map[uint64]*common.CloseWrapper[*dto.Message]{},
		Udp:     udp,
	}
}

func (lrcp *LRCP) SendToSession(msg *dto.Message) {
	log.Println("3")

	session, ok := lrcp.Session[msg.SessionId]
	log.Println("4")

	// If session does not exist then create one
	if !ok {
		// If first message is not connect: send /close/SESSION/ and stop.
		if msg.MsgType != "connect" {
			lrcp.sendClose(msg)
			return
		}
		log.Println("CREATING NEW SESSION")
		session = common.NewCloseWrapper[*dto.Message](CHANNEL_BUFFER)
		lrcp.Session[msg.SessionId] = session

		go lrcp.handleSession(session, msg.SessionId, msg.Addr)
	}
	log.Println("77")

	// Send to the session
	session.Send(msg)
}

func (lrcp *LRCP) handleSession(session *common.CloseWrapper[*dto.Message], sessionId uint64, addr *net.UDPAddr) {
	defer func() {
		log.Println("CLOSINGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG")
		session.Close()
	}()

	var lengthReceived uint64 = 0

	timer := time.NewTimer(SESSION_EXPIRY_TIMEOUT)

	pr, pw := io.Pipe()

	ack := make(chan *dto.Message, CHANNEL_BUFFER)
	go func() {
		err := lrcp.reverse(pr, sessionId, addr, ack, session.Done())
		log.Println("3333 CLOSINGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG")

		if err != nil {
			log.Println("2222 CLOSINGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG")
			session.Close()
		}
	}()

	var msg *dto.Message
	for {

		select {
		case _msg := <-session.Receive():

			log.Println("5")
			msg = _msg
			log.Println("6")
			timer.Reset(SESSION_EXPIRY_TIMEOUT)
		case <-timer.C:
			return
		case <-session.Done():
			return
		}

		if msg.MsgType == "close" {
			lrcp.sendClose(msg)
			break
		}

		if msg.MsgType == "connect" {
			lrcp.sendAck(msg, 0)
		}

		if msg.MsgType == "data" {
			lengthReceived = lrcp.handleData(msg, lengthReceived, pw)
		}

		// log.Println(msg.MsgType, msg.SessionId, string(msg.Buf))
		if msg.MsgType == "ack" {
			select {
			case ack <- msg:
			default:
				log.Println("DROPPED")
			}
		}

	}
}

func (lrcp *LRCP) sendClose(msg *dto.Message) {
	closeMsg := fmt.Sprintf("/close/%d/", msg.SessionId)
	lrcp.Udp.WriteToUDP([]byte(closeMsg), msg.Addr)
}

func (lrcp *LRCP) sendAck(msg *dto.Message, lengthReceived uint64) {
	ackMsg := fmt.Sprintf("/ack/%d/%d/", msg.SessionId, lengthReceived)
	lrcp.Udp.WriteToUDP([]byte(ackMsg), msg.Addr)
}

func (lrcp *LRCP) handleData(msg *dto.Message, lengthReceived uint64, pw *io.PipeWriter) uint64 {

	str := string(msg.Buf)

	nextPosStr, payload, found := strings.Cut(str, "/")
	if !found {
		return lengthReceived
	}

	nextPos, err := strconv.Atoi(nextPosStr)
	if err != nil {
		return lengthReceived
	}

	if lengthReceived < uint64(nextPos) {
		lrcp.sendAck(msg, lengthReceived)
		return lengthReceived
	}
	if lengthReceived > uint64(nextPos) {
		return lengthReceived
	}

	// Remove the last '/'
	payload = payload[:len(payload)-1]

	payload = strings.ReplaceAll(payload, `\/`, `/`)
	payload = strings.ReplaceAll(payload, `\\`, `\`)

	totalLength := lengthReceived + uint64(len(payload))
	lrcp.sendAck(msg, totalLength)

	pw.Write([]byte(payload))

	return totalLength
}
