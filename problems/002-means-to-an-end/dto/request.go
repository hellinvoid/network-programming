package dto

import (
	"encoding/binary"
	"io"
)
// Request length defined in question
const MESSAGE_SIZE = 9

type Request struct {
	MessageType string
	Timestamp   int32
	Price       int32
	MinTime     int32
	MaxTime     int32
}

// Function to decode the bytes into request
func NewRequest(r io.Reader) (*Request, error) {
	buffer := make([]byte, MESSAGE_SIZE)
	_, err := io.ReadFull(r, buffer)
	if err != nil {
		return nil, err
	}

	req := &Request{}
	req.MessageType = string(buffer[0])

	// BigEndian conversion
	first := binary.BigEndian.Uint32(buffer[1:])
	second := binary.BigEndian.Uint32(buffer[5:])

	if req.MessageType == "I" {
		req.Timestamp = int32(first)
		req.Price = int32(second)
	} else {
		req.MinTime = int32(first)
		req.MaxTime = int32(second)
	}

	return req, nil
}
