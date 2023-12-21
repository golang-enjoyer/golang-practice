package main

import (
	"fmt"
	"sync"
	"time"
)

var ErrNoWorkers = fmt.Errorf("attempting to create worker pool with less than 1 worker")
var ErrNegativeChannelSize = fmt.Errorf("attempting to create worker pool with a negative channel size")

func NewSimplePool(numWorkers int, channelSize int) (Pool, error) {
	if numWorkers <= 0 {
		return nil, ErrNoWorkers
	}
	if channelSize < 0 {
		return nil, ErrNegativeChannelSize
	}

	tasks := make(chan Task, channelSize)

	return &WorkerPool{
		numWorkers: numWorkers,
		tasks:      tasks,

		start: sync.Once{},
		end:   sync.Once{},

		quit: make(chan struct{}),
	}, nil
}

func (p *WorkerPool) Start() {
	p.start.Do(func() {
		p.startWorkers()
	})
}

func (p *WorkerPool) Stop() {
	p.end.Do(func() {
		close(p.quit)
	})
}

func (p *WorkerPool) AddTask(t Task) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

func (p *WorkerPool) startWorkers() {
	for i := 0; i < p.numWorkers; i++ {
		go func(workerNum int) {

			for {
				select {
				case <-p.quit:
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}

					if err := task.Execute(); err != nil {
						task.HandleError(err)
					}
				}
			}
		}(i)
	}
}

type MyTask struct {
	ID int
}

func (t *MyTask) Execute() error {
	fmt.Printf("Task %d is being executed\n", t.ID)
	time.Sleep(time.Second)
	fmt.Printf("Task %d executed successfully\n", t.ID)
	return nil
}

func (t *MyTask) HandleError(err error) {
	fmt.Printf("Task %d failed: %v\n", t.ID, err)
}

func main() {
	numWorkers := 3
	channelSize := 5

	pool, err := NewSimplePool(numWorkers, channelSize)
	if err != nil {
		fmt.Printf("Error creating pool: %v\n", err)
		return
	}

	pool.Start()

	for i := 0; i < 10; i++ {
		task := &MyTask{ID: i}
		pool.AddTask(task)
	}

	time.Sleep(5 * time.Second)

	pool.Stop()
}
