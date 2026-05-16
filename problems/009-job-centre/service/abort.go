package service

import (
	"fmt"
	"log"
	"net"

	"github.com/hellinvoid/network-programming/problems/009-job-centre/dto"
)

func (jc *JobCentre) HandleAbort(m *dto.Message, conn net.Conn, allocated []*QueueEntry) []*QueueEntry {
	if m.Id == nil {
		jc.SendMessage(conn, dto.ERROR_RESPONSE)
		return allocated
	}

	allocatedToMe := false
	// log.Println(allocated)
	for index, entry := range allocated {
		if entry.Id == *m.Id {
			allocated = append(allocated[:index], allocated[index+1:]...)
			allocatedToMe = true
			break
		}
	}
	log.Println(allocatedToMe)
	q, ok := jc.QueueFromId[*m.Id]

	if ok && allocatedToMe {
		if q.Abort(*m.Id, conn.RemoteAddr().String(), false) {
			jc.SendMessage(conn, fmt.Sprintf(dto.ABORT_RESPONSE, "ok"))
		} else {
			jc.SendMessage(conn, fmt.Sprintf(dto.ABORT_RESPONSE, "no-job"))
		}
	} else if ok && q.IsAllocated(*m.Id) {
		jc.SendMessage(conn, dto.ERROR_RESPONSE)
	} else {
		jc.SendMessage(conn, fmt.Sprintf(dto.ABORT_RESPONSE, "no-job"))
	}
	return allocated
}

func (jc *JobCentre) AbortAll(conn net.Conn, allocated []*QueueEntry) {
	addr := conn.RemoteAddr().String()

	for _, entry := range allocated {
		q, ok := jc.QueueFromId[entry.Id]
		if !ok {
			continue
		}
		q.Abort(entry.Id, addr, true)
	}
}
