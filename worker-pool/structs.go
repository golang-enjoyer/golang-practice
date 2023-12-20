package main

import "sync"

type WorkerPool struct {
	workersAmount int
	tasks         chan Task

	start sync.Once
	end   sync.Once
	quit  chan struct{}
}
