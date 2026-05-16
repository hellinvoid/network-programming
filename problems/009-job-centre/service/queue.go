package service

import (
	"log"
	"sync"
)

const (
	CHANNEL_BUFFER_SIZE = 1000
)

type QueueEntry struct {
	Id  int
	Job any
	Pri int
}

type Queue struct {
	name string

	mu        sync.Mutex
	entries   []*QueueEntry
	allocated map[string][]*QueueEntry
}

func NewQueue(name string) *Queue {
	return &Queue{
		name:      name,
		mu:        sync.Mutex{},
		entries:   make([]*QueueEntry, 0),
		allocated: map[string][]*QueueEntry{},
	}
}

func (q *Queue) Add(qe ...*QueueEntry) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.entries = append(q.entries, qe...)
}

func (q *Queue) Delete(id int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for index, entry := range q.entries {
		if entry.Id == id {
			q.entries = append(q.entries[:index], q.entries[index+1:]...)
			return true
		}
	}

	for key, qeArr := range q.allocated {
		for index, entry := range qeArr {
			if entry.Id == id {
				q.allocated[key] = append(qeArr[:index], qeArr[index+1:]...)
				return true
			}
		}
	}

	return false
}

func (q *Queue) Abort(id int, addr string, all bool) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	qeArr, ok := q.allocated[addr]
	if !ok {
		return false
	}
	if all {
		q.entries = append(q.entries, qeArr...)

		delete(q.allocated, addr)
		return true
	}

	for index, entry := range qeArr {
		if entry.Id == id {
			q.allocated[addr] = append(qeArr[:index], qeArr[index+1:]...)
			q.entries = append(q.entries, entry)

			return true
		}
	}
	return false
}

func (q *Queue) IsAllocated(id int) bool {
	for _, qeArr := range q.allocated {
		for _, entry := range qeArr {
			if entry.Id == id {
				return true
			}
		}
	}

	return false
}

func (q *Queue) Get(addr string) *QueueEntry {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.entries) == 0 {
		log.Println("entries empty")
		return nil
	}
	max := 0
	for index, entry := range q.entries {
		if entry.Pri > q.entries[max].Pri {
			max = index
		}
	}
	entry := q.entries[max]

	return entry
}

func (q *Queue) Allocate(entry *QueueEntry, addr string) bool {
	if !q.Delete(entry.Id) {
		return false
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	_, ok := q.allocated[addr]
	if !ok {
		q.allocated[addr] = make([]*QueueEntry, 0)
	}
	q.allocated[addr] = append(q.allocated[addr], entry)
	return true
}
