package dto

import (
	"encoding/binary"
	"io"
)

// Function to encode and send the integer response
func SendResponse(res int32, w io.Writer) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(res))
	
	_, err := w.Write(buf)
	return err
}
