package dto

import (
	"encoding/json"
)

const (
	ERROR_RESPONSE       = "{\"status\":\"error\",\"error\":\"You did something wrong.\"}\n"
	PUT_RESPONSE         = "{\"status\":\"ok\",\"id\":%d}\n"
	DELETE_RESPONSE      = "{\"status\":\"%s\"}\n"
	ABORT_RESPONSE       = "{\"status\":\"%s\"}\n"
	GET_RESPONSE_SUCCESS = "{\"status\":\"ok\",\"id\":%d,\"job\":%s,\"pri\":%d,\"queue\":\"%s\"}\n"
	GET_RESPONSE_FAILURE = "{\"status\":\"no-job\"}\n"
)

type Message struct {
	// all
	Request string `json:"request"`

	// put
	Queue *string `json:"queue"`
	Job   any     `json:"job"`
	Pri   *int    `json:"pri"`

	// get
	Queues []string `json:"queues"`
	Wait   bool     `json:"wait"`

	// delete, abort
	Id *int `json:"id"`
}

func NewMessage(d *json.Decoder) (*Message, error) {
	m := &Message{}
	err := d.Decode(m)
	return m, err
}
