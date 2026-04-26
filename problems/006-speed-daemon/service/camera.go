package service

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

// Camera client
type Camera struct {
	road      Road
	mile      Mile
	limit     Limit
	conn      net.Conn
	ts        *TicketSystem
	heartbeat bool
}

func NewCamera(conn net.Conn, ts *TicketSystem, heartbeat bool, r *bufio.Reader) (*Camera, error) {

	buf := make([]byte, SIZE_u16*3)

	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	cm := &Camera{
		road:      Road(binary.BigEndian.Uint16(buf[0:2])),
		mile:      Mile(binary.BigEndian.Uint16(buf[2:4])),
		limit:     Limit(binary.BigEndian.Uint16(buf[4:6])),
		conn:      conn,
		ts:        ts,
		heartbeat: heartbeat,
	}

	return cm, nil

}

func (cm *Camera) HandleConn(r *bufio.Reader) error {

	for {
		code, err := r.ReadByte()
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		log.Println(getCode(code), code)

		if code == WANT_HEARTBEAT && !cm.heartbeat {
			SendHeartbeat(cm.conn, r)
			cm.heartbeat = true
		} else if code == PLATE {
			cm.DecodePlate(r)
		} else {
			return errors.New("Invalid message")
		}
	}

}

// Decode plate message
func (cm *Camera) DecodePlate(r *bufio.Reader) {
	b, err := r.ReadByte()
	if err != nil {
		return
	}
	plateLen := int(b)

	buf := make([]byte, plateLen+SIZE_u32)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return
	}
	// Pass on the plate message to ticket system for further processing
	plate := Plate(buf[:plateLen])
	timestamp := Timestamp(binary.BigEndian.Uint32(buf[plateLen:]))
	sp := &SpeedProps{
		plate:      plate,
		road:       cm.road,
		mile:       cm.mile,
		timestamp:  timestamp,
		speedLimit: cm.limit,
	}
	cm.ts.spch <- sp
}
