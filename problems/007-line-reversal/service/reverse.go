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

func (lrcp *LRCP) reverse(pr *io.PipeReader, sessionId uint64, addr *net.UDPAddr, ack chan *dto.Message, end chan any) error {
	r := bufio.NewReader(pr)

	var fullPayLoad strings.Builder
	maxLengthAck := 0

	posMap := map[int]int{}
	pos := 0
	buf := make([]byte, 1000)
	for {

		posStr, err := r.ReadString('/')
		if err != nil {
			return err
		}
		pos, err = strconv.Atoi(posStr[:len(posStr)-1])
		if err != nil {
			return err
		}

		n, err := r.Read(buf)
		if err != nil {
			return err
		}

		payload := rev(string(buf[:n]))
		fullPayLoad.WriteString(payload)
		cumulativeLenght := fullPayLoad.Len()
		log.Println(cumulativeLenght)
		posMap[cumulativeLenght] = pos

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
				// '/' must be at end so removing it
				str := string(ackMsg.Buf[:len(ackMsg.Buf)-1])

				// Get length
				length, err := strconv.Atoi(str)
				if err != nil {
					return err
				}

				// If length is as expected stop retransmission
				if length == len(fullPayLoad.String()) {
					retransmit = false
					ticker.Stop()
				} else if length > maxLengthAck && length < len(fullPayLoad.String()) {
					reqPos := posMap[length]
					err := sendReverse(fullPayLoad.String()[length+1:], lrcp.Udp, sessionId, addr, reqPos)
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

func sendReverse(payload string, udp *net.UDPConn, sessionId uint64, addr *net.UDPAddr, pos int) error {
	payload = strings.ReplaceAll(payload, `\`, `\\`)
	payload = strings.ReplaceAll(payload, `/`, `\/`)

	dataMsg := fmt.Sprintf("/data/%d/%d/%s/", sessionId, pos, payload)
	log.Println("Sending ", dataMsg)

	_, err := udp.WriteTo([]byte(dataMsg), addr)
	return err
}

func rev(strs string) string {

	var last strings.Builder
	for str := range strings.SplitSeq(strs, "\n") {
		r := []rune(str)
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}

		last.WriteString(string(r) + "\n")
	}

	return last.String()[:last.Len()-1]
}
