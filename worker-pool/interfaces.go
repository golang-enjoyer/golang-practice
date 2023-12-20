package main

type Pool interface {
	Start()
	Close()
	AddTask(Task)
}

type Task interface {
	Execute()
	HandleError(error)
}
