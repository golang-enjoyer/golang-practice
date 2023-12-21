package main

type Pool interface {
	Start()
	Stop()
	AddTask(Task)
}

type Task interface {
	Execute() error
	HandleError(error)
}
