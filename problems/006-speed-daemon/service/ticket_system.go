package service

import (
	"bufio"
	"log"
	"net"
	"sync"
)

// Message codes
const (
	// Error code
	ERROR = 0x10

	// Codes related to camera
	I_AM_CAMERA    = 0x80
	PLATE          = 0x20
	WANT_HEARTBEAT = 0x40
	HEARTBEAT      = 0x41

	// Codes related to dispatcher
	I_AM_DISPATCHER = 0x81
	TICKET          = 0x21
)

// Size bits -> bytes mapping
const (
	SIZE_u16 = 2
	SIZE_u32 = 4
)

const DAY_DIV = 86400

// Type specifiers
type (
	Plate     string
	Road      uint16
	Mile      uint16
	Limit     uint16
	Timestamp uint32
	Day       uint32
)

// Speed props (data to pass along the channels)
type SpeedProps struct {
	road       Road
	plate      Plate
	mile       Mile
	timestamp  Timestamp
	speedLimit Limit
}

type TicketSystem struct {
	spch                chan *SpeedProps
	hasTicketOnDay      map[string]any
	roadPlateSpeedProps map[string][]*SpeedProps

	mu         sync.Mutex
	dispatcher map[Road]chan net.Conn
}

func NewTicketSystem() *TicketSystem {
	return &TicketSystem{
		spch:                make(chan *SpeedProps),
		hasTicketOnDay:      map[string]any{},
		roadPlateSpeedProps: map[string][]*SpeedProps{},

		mu:         sync.Mutex{},
		dispatcher: map[Road]chan net.Conn{},
	}
}

// Function to handle and route the incoming client asper camera or ticket dispatcher
func (ts *TicketSystem) HandleConn(conn net.Conn) {
	defer conn.Close()
	heartbeat := false
	// Identify camera or dispatcher
	r := bufio.NewReader(conn)
	for {

		code, err := r.ReadByte()
		if err != nil {
			return
		}
		log.Println(getCode(code))

		switch code {

		case I_AM_CAMERA:
			cm, err := NewCamera(conn, ts, heartbeat, r)
			if err != nil {
				return
			}
			ts.getDispatcher(cm.road)
			// Let the camera handle connection
			err = cm.HandleConn(r)
			if err != nil {
				SendError(conn)
			}
			return
		case I_AM_DISPATCHER:
			td, err := NewTicketDispatcher(conn, heartbeat, r)
			if err != nil {
				return
			}
			for _, road := range td.roads {
				log.Println("Conn: ", conn)
				ch := ts.getDispatcher(road)
				ch <- conn
			}
			err = td.HandleConn(r)
			if err != nil {
				SendError(conn)
			}
			return
		case WANT_HEARTBEAT:
			if heartbeat {
				SendError(conn)
				return
			}
			heartbeat = true
			SendHeartbeat(conn, r)
		default:
			SendError(conn)
		}
	}

}

// Function to send error on violation of protocol
func SendError(conn net.Conn) {
	// If we send an error means we close the connection
	defer conn.Close()

	errorMsg := "Invalid"

	msg := make([]byte, 0)
	msg = append(msg, ERROR, byte(len(errorMsg)))
	msg = append(msg, []byte(errorMsg)...)

	conn.Write(msg)
}

func getCode(code uint8) string {
	switch code {
	case ERROR:
		return "ERROR"

	// Codes related to camera
	case I_AM_CAMERA:
		return "I_AM_CAMERA"
	case PLATE:
		return "PLATE"
	case WANT_HEARTBEAT:
		return "WANT_HEARTBEAT"
	case HEARTBEAT:
		return "HEARTBEAT"

	// Codes related to dispatcher
	case I_AM_DISPATCHER:
		return "I_AM_DISPATCHER"
	case TICKET:
		return "TICKET"
	default:
	}
	return "INVALID"
}

func (ts *TicketSystem) getDispatcher(road Road) chan net.Conn {
	defer ts.mu.Unlock()
	ts.mu.Lock()

	ch, ok := ts.dispatcher[road]
	// Initialize the dispatcher if not there
	if !ok {
		ch = make(chan net.Conn, 10)
		ts.dispatcher[road] = ch
	}
	return ch
}
