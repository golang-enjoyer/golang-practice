package main

import "sync"

type WorkerPool struct {
	numWorkers int
	tasks      chan Task

	start sync.Once
	quit  chan struct{}
	end   sync.Once
}
