package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hellinvoid/network-programming/problems/009-job-centre/dto"
)

type GetResponse struct {
	qname string
	qe    *QueueEntry
}

func (jc *JobCentre) HandleGet(m *dto.Message, conn net.Conn, allocated []*QueueEntry) []*QueueEntry {
	log.Println("Trying 19")
	if len(m.Queues) == 0 {
		log.Println("Q seems to be empty")
		jc.SendMessage(conn, dto.ERROR_RESPONSE)
		return allocated
	}

	addr := conn.RemoteAddr().String()
	grs := make([]*GetResponse, 0)

	grCh := make(chan *GetResponse)

	for _, qname := range m.Queues {
		jc.mu.Lock()
		q, ok := jc.AllQueues[qname]
		if !ok {
			q = NewQueue(qname)
			jc.AllQueues[qname] = q
		}
		jc.mu.Unlock()

		log.Println("Trying to get a job")
		qe1 := q.Get(addr)
		if qe1 != nil {
			grs = append(grs, &GetResponse{qname: qname, qe: qe1})
			continue
		}
		if m.Wait {
			go func(q *Queue, grCh chan *GetResponse) {
				ticker := time.NewTicker(time.Second)
				for range ticker.C {
					res := q.Get(addr)
					if res != nil {
						grCh <- &GetResponse{
							qname: q.name,
							qe:    res,
						}
						return
					}
				}

			}(q, grCh)
		}
	}

	for len(grs) > 0 {
		max := 0
		for index, gr := range grs {
			if gr.qe.Pri > grs[max].qe.Pri {
				max = index
			}
		}
		gr := grs[max]
		grs = append(grs[:max], grs[max+1:]...)

		q := jc.AllQueues[gr.qname]
		if !q.Allocate(gr.qe, addr) {
			continue
		}

		allocated = append(allocated, gr.qe)
		job, _ := json.Marshal(gr.qe.Job)
		msg := fmt.Sprintf(dto.GET_RESPONSE_SUCCESS, gr.qe.Id, job, gr.qe.Pri, gr.qname)
		jc.SendMessage(conn, msg)
		return allocated
	}

	if !m.Wait {
		jc.SendMessage(conn, dto.GET_RESPONSE_FAILURE)
		return allocated
	}

	for gr := range grCh {
		q := jc.AllQueues[gr.qname]
		if !q.Allocate(gr.qe, addr) {
			continue
		}
		allocated = append(allocated, gr.qe)
		job, _ := json.Marshal(gr.qe.Job)
		msg := fmt.Sprintf(dto.GET_RESPONSE_SUCCESS, gr.qe.Id, job, gr.qe.Pri, gr.qname)
		jc.SendMessage(conn, msg)

		return allocated

	}
	return allocated
}
