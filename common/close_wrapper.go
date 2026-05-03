package common

import (
	"log"
	"sync"
)

type CloseWrapper[T any] struct {
	ch   chan T
	done chan any
	once sync.Once
}

func NewCloseWrapper[T any](size int) *CloseWrapper[T] {
	return &CloseWrapper[T]{
		ch:   make(chan T, size),
		done: make(chan any),
	}
}

func (cw *CloseWrapper[T]) Send(val T) {
	select {
	case cw.ch <- val:
		log.Println("SENT")
	case <-cw.done:
		log.Println("DROPPED")
	}
}

func (cw *CloseWrapper[T]) Receive() <-chan T {
	return cw.ch
}

func (cw *CloseWrapper[T]) Close() {
	cw.once.Do(func() {
		close(cw.done)
	})
}

func (cw *CloseWrapper[T]) Done() <-chan any {
	return cw.done
}
