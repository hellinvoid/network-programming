package service

import (
	"net"
	"sync"
)

type JobCentre struct {
	mu sync.Mutex

	IdCounter   int
	AllQueues   map[string]*Queue
	QueueFromId map[int]*Queue
}

func NewJobCentre() *JobCentre {
	return &JobCentre{
		mu:          sync.Mutex{},
		IdCounter:   1,
		AllQueues:   map[string]*Queue{},
		QueueFromId: map[int]*Queue{},
	}
}

func (jc *JobCentre) SendMessage(conn net.Conn, msg string) {
	conn.Write([]byte(msg))
}
