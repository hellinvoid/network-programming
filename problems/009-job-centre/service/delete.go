package service

import (
	"fmt"
	"net"

	"github.com/hellinvoid/network-programming/problems/009-job-centre/dto"
)

func (jc *JobCentre) HandleDelete(m *dto.Message, conn net.Conn) {
	if m.Id == nil {
		jc.SendMessage(conn, dto.ERROR_RESPONSE)
		return
	}

	q, ok := jc.QueueFromId[*m.Id]
	if !ok {
		jc.SendMessage(conn, fmt.Sprintf(dto.DELETE_RESPONSE, "no-job"))
		return
	}

	if !q.Delete(*m.Id) {
		jc.SendMessage(conn, fmt.Sprintf(dto.DELETE_RESPONSE, "no-job"))
		return
	}
	delete(jc.QueueFromId, *m.Id)
	jc.SendMessage(conn, fmt.Sprintf(dto.DELETE_RESPONSE, "ok"))
}
