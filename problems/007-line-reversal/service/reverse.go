package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hellinvoid/network-programming/problems/007-line-reversal/dto"
)

func (lrcp *LRCP) reverse(pr *io.PipeReader, sessionId uint64, addr *net.UDPAddr, ack chan *dto.Message, end <-chan any) error {
	r := bufio.NewReader(pr)

	var fullPayLoad strings.Builder
	maxLengthAck := 0

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return err
		}

		pos := fullPayLoad.Len()
		payload := rev(string(line)) + "\n"
		fullPayLoad.WriteString(payload)
		log.Println("TOTAL :", fullPayLoad.Len())

		sendReverse(payload, lrcp.Udp, sessionId, addr, pos)
		// Start listening for ack else retransmit
		ticker := time.NewTicker(RETRANSMISSION_TIMEOUT)

		var ackMsg *dto.Message
		var retransmit bool = true

		for retransmit {
			select {
			case <-ticker.C:
				// Retransmit previous message
				err := sendReverse(payload, lrcp.Udp, sessionId, addr, pos)
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
				log.Println("Received :", length)
				if length == fullPayLoad.Len() {
					retransmit = false
					ticker.Stop()

				} else if length >= maxLengthAck && length < fullPayLoad.Len() {

					err := sendReverse(fullPayLoad.String()[length:], lrcp.Udp, sessionId, addr, length)
					if err != nil {
						return err
					}

				} else if length > fullPayLoad.Len() {
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

func sendReverse(payload string, udp *net.UDPConn, sessionId uint64, addr *net.UDPAddr, pos int) error {
	payload = strings.ReplaceAll(payload, `\`, `\\`)
	payload = strings.ReplaceAll(payload, `/`, `\/`)

	dataMsg := fmt.Sprintf("/data/%d/%d/%s/", sessionId, pos, payload)
	// log.Println("Sending ", dataMsg)

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
