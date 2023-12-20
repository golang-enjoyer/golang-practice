package main

import (
	"sync"
)

func NewSimplePool(workerAmount int, chanSize int) (Pool, error) {
	tasks := make(chan Task, chanSize)

	return &WorkerPool{
		workersAmount: workerAmount,
		tasks:         tasks,
		start:         sync.Once{},
		end:           sync.Once{},
		quit:          make(chan struct{}),
	}, nil
}

func (p *WorkerPool) Start() {
	p.start.Do(func() {
		p.start()
	})
}

func (p *WorkerPool) Stop() {
	p.end.Do(func() {
		close(p.quit)
	})
}

func (p *WorkerPool) AddWork(t Task) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}
