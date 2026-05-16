package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/hellinvoid/network-programming/problems/009-job-centre/dto"
)

func (jc *JobCentre) HandlePut(m *dto.Message, conn net.Conn) {
	if !isValidPut(m) {
		log.Println("Invalid Message")
		jc.SendMessage(conn, dto.ERROR_RESPONSE)
		return
	}

	qe := &QueueEntry{
		Job: m.Job,
		Pri: *m.Pri,
	}

	jc.mu.Lock()

	q, ok := jc.AllQueues[*m.Queue]
	if !ok {
		q = NewQueue(*m.Queue)
		jc.AllQueues[*m.Queue] = q
	}

	qe.Id = jc.IdCounter
	jc.QueueFromId[qe.Id] = q
	
	jc.IdCounter++
	
	jc.mu.Unlock()

	q.Add(qe)

	jc.SendMessage(conn, fmt.Sprintf(dto.PUT_RESPONSE, qe.Id))
}

func isValidPut(m *dto.Message) bool {
	if m.Queue == nil || m.Job == nil || m.Pri == nil {
		return false
	}

	_, err := json.Marshal(m.Job)

	return err == nil && *m.Pri >= 0
}
