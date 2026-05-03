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
	session, ok := lrcp.Session[msg.SessionId]
	// If session does not exist then create one
	if !ok {
		// If first message is not connect: send /close/SESSION/ and stop.
		if msg.MsgType != "connect" {
			lrcp.sendClose(msg)
			return
		}
		log.Println("CREATING NEW SESSION")
		session = common.NewCloseWrapper[*dto.Message]()
		lrcp.Session[msg.SessionId] = session

		go lrcp.handleSession(session, msg.SessionId, msg.Addr)
	}
	// Send to the session
	session.Send(msg)
}

func (lrcp *LRCP) handleSession(session *common.CloseWrapper[*dto.Message], sessionId uint64, addr *net.UDPAddr) {
	end := make(chan any)
	defer func() {
		session.Close()
		close(end)
	}()

	var lengthReceived uint64 = 0
	var pos uint64 = 0

	timer := time.NewTimer(SESSION_EXPIRY_TIMEOUT)

	pr, pw := io.Pipe()

	ack := make(chan *dto.Message)
	go func() {
		err := lrcp.reverse(pr, sessionId, addr, ack, end)
		if err != nil {
			session.Close()
		}
	}()

	var msg *dto.Message
	for {

		select {
		case _msg, ok := <-session.Receive():
			if !ok {
				return
			}
			msg = _msg
			timer.Reset(SESSION_EXPIRY_TIMEOUT)
		case <-timer.C:
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
			pos, lengthReceived = lrcp.handleData(msg, pos, lengthReceived, pw)
		}

		if msg.MsgType == "ack" {
			ack <- msg
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

func (lrcp *LRCP) handleData(msg *dto.Message, pos uint64, lengthReceived uint64, pw *io.PipeWriter) (uint64, uint64) {

	str := string(msg.Buf)

	nextPosStr, payload, found := strings.Cut(str, "/")
	if !found {
		return pos, lengthReceived
	}

	nextPos, err := strconv.Atoi(nextPosStr)
	if err != nil {
		return pos, lengthReceived
	}

	if pos < uint64(nextPos) {
		lrcp.sendAck(msg, lengthReceived)
		return pos, lengthReceived
	}
	if pos > uint64(nextPos) {
		return pos, lengthReceived
	}

	// Remove the last '/'
	payload = payload[:len(payload)-1]

	payload = strings.ReplaceAll(payload, `\/`, `/`)
	payload = strings.ReplaceAll(payload, `\\`, `\`)

	totalLength := lengthReceived + uint64(len(payload))
	lrcp.sendAck(msg, totalLength)

	fmt.Fprint(pw, pos)
	fmt.Fprint(pw, "/")
	fmt.Fprint(pw, payload)

	return uint64(nextPos) + 1, totalLength
}
