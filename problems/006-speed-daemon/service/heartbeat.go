package service

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"time"
)

// Function to send heartbeat to the client
func SendHeartbeat(conn net.Conn, r *bufio.Reader) {
	buf := make([]byte, SIZE_u32)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return
	}
	interval := binary.BigEndian.Uint32(buf)
	if interval == 0 {
		return
	}
	intervalInSecs := float32(interval) / 10
	ticker := time.NewTicker(time.Duration(intervalInSecs) * time.Second)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			_, err := conn.Write([]byte{HEARTBEAT})
			if err != nil {
				return
			}
		}
	}()
}
