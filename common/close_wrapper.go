package common

import (
	"sync"
)

type CloseWrapper[T any] struct {
	ch       chan T
	mu       sync.Mutex
	isClosed bool
}

func NewCloseWrapper[T any]() *CloseWrapper[T] {
	return &CloseWrapper[T]{
		ch:       make(chan T),
		mu:       sync.Mutex{},
		isClosed: false,
	}
}

func (cw *CloseWrapper[T]) Send(val T) {
	cw.mu.Lock()
	if !cw.isClosed {
		cw.ch <- val
	}
	cw.mu.Unlock()
}

func (cw *CloseWrapper[T]) Receive() <-chan T {
	return cw.ch
}

func (cw *CloseWrapper[T]) Close() {
	cw.mu.Lock()
	close(cw.ch)
	cw.isClosed = true
	cw.mu.Unlock()
}
