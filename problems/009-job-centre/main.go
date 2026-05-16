package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"

	"github.com/hellinvoid/network-programming/problems/009-job-centre/dto"
	"github.com/hellinvoid/network-programming/problems/009-job-centre/service"
)

const (
	PROTOCOL = "tcp"
	ADDRESS  = ":6969"
)

func main() {
	listener, err := net.Listen(PROTOCOL, ADDRESS)
	if err != nil {
		panic(err)
	}
	log.Println("Server listening on", listener.Addr().String())

	jc := service.NewJobCentre()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		log.Println(conn.RemoteAddr().String(), " Connected")
		go handleConn(conn, jc)
	}
}

func handleConn(conn net.Conn, jc *service.JobCentre) {
	defer conn.Close()

	d := json.NewDecoder(conn)

	allocated := make([]*service.QueueEntry, 0)
	defer jc.AbortAll(conn, allocated)

	for {
		// conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		m, err := dto.NewMessage(d)
		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				return
			}
			log.Println("Error decoding", err)
			jc.SendMessage(conn, dto.ERROR_RESPONSE)
			continue
		}
		switch m.Request {
		case "put":
			jc.HandlePut(m, conn)
		case "get":
			log.Println("Trying Get")
			allocated = jc.HandleGet(m, conn, allocated)
		case "delete":
			jc.HandleDelete(m, conn)
		case "abort":
			allocated = jc.HandleAbort(m, conn, allocated)
		default:
			jc.SendMessage(conn, dto.ERROR_RESPONSE)
		}

	}

}
