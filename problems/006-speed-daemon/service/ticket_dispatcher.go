package service

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

// TicketDispatcher client
type TicketDispatcher struct {
	conn      net.Conn
	roads     []Road
	heartbeat bool
}

func NewTicketDispatcher(conn net.Conn, heartbeat bool, r *bufio.Reader) (*TicketDispatcher, error) {

	n, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, SIZE_u16*int(n))

	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	td := &TicketDispatcher{
		conn:      conn,
		roads:     make([]Road, 0),
		heartbeat: heartbeat,
	}
	for en := SIZE_u16; en <= SIZE_u16*int(n); en += SIZE_u16 {
		td.roads = append(td.roads, Road(binary.BigEndian.Uint16(buf[en-SIZE_u16:en])))
	}
	log.Println(td.roads)
	return td, err
}

func (td *TicketDispatcher) HandleConn(r *bufio.Reader) error {

	for {
		code, err := r.ReadByte()
		if err != nil {
			return nil
		}

		if code == WANT_HEARTBEAT && !td.heartbeat {
			SendHeartbeat(td.conn, r)
			td.heartbeat = true
		} else {
			return errors.New("Invalid message")
		}
	}

}

type Ticket struct {
	plate      Plate
	road       Road
	mile1      Mile
	timestamp1 Timestamp
	mile2      Mile
	timestamp2 Timestamp
	speed      Limit
}

func (t *Ticket) Encode() []byte {
	buf := make([]byte, 0)

	buf = append(buf, TICKET, byte(len(t.plate)))
	buf = append(buf, []byte(t.plate)...)

	buf = binary.BigEndian.AppendUint16(buf, uint16(t.mile1))
	buf = binary.BigEndian.AppendUint32(buf, uint32(t.timestamp1))

	buf = binary.BigEndian.AppendUint16(buf, uint16(t.mile2))
	buf = binary.BigEndian.AppendUint32(buf, uint32(t.timestamp2))

	buf = binary.BigEndian.AppendUint16(buf, uint16(t.speed))

	return buf
}
