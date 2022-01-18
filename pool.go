package main

import (
	"sync"
	"sync/atomic"
)

type pool struct {
	status atomic.Value
	queue  chan int
	wg     *sync.WaitGroup
}

func EzPool(size int) *pool {
	if size <= 0 {
		size = 1
	}
	p := &pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
	p.status.Store(true)
	return p
}

func (p *pool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	p.wg.Add(delta)
}

func (p *pool) Done() {
	<-p.queue
	p.wg.Done()
}

func (p *pool) Wait() {
	if !p.status.Load().(bool) {
		return
	}
	p.wg.Wait()
}

func (p *pool) Close() {
	p.status.Store(false)
}