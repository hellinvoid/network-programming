package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hellinvoid/network-programming/problems/007-line-reversal/dto"
)

func (lrcp *LRCP) reverse(pr *io.PipeReader, sessionId uint64, addr *net.UDPAddr, ack chan *dto.Message, end chan any) error {
	r := bufio.NewReader(pr)

	var fullPayLoad strings.Builder
	maxLengthAck := 0

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return err
		}

		payload := rev(string(line)) + "\n"
		fullPayLoad.WriteString(payload)

		sendReverse(payload, lrcp.Udp, sessionId, addr)
		// Start listening for ack else retransmit
		ticker := time.NewTicker(RETRANSMISSION_TIMEOUT)

		var ackMsg *dto.Message
		var retransmit bool = true

		for retransmit {
			select {
			case <-ticker.C:
				// Retransmit previous message
				err := sendReverse(payload, lrcp.Udp, sessionId, addr)
				if err != nil {
					return err
				}
			case ackMsg = <-ack:
				// Decode ack message
				str := string(ackMsg.Buf[:len(ackMsg.Buf)-1])

				length, err := strconv.Atoi(str)
				if err != nil {
					return err 
				}
				if length == len(fullPayLoad.String()) {
					retransmit = false
					ticker.Stop()

				} else if length > maxLengthAck && length < len(fullPayLoad.String()) {
					err := sendReverse(fullPayLoad.String()[length+1:], lrcp.Udp, sessionId, addr)
					if err != nil {
						return err
					}

				} else if length > len(fullPayLoad.String()) {
					return errors.New("Misbehaving")
				}
				if length > maxLengthAck {
					maxLengthAck = length
				}
			case <-end:
				return nil
			}
		}

	}
}

func sendReverse(payload string, udp *net.UDPConn, sessionId uint64, addr *net.UDPAddr) error {
	payload = strings.ReplaceAll(payload, `\`, `\\`)
	payload = strings.ReplaceAll(payload, `/`, `\/`)

	dataMsg := fmt.Sprintf("/data/%d/0/%s/", sessionId, payload)

	_, err := udp.WriteTo([]byte(dataMsg), addr)
	return err
}

func rev(str string) string {
	r := []rune(str)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
