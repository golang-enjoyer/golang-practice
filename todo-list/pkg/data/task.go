package data

import (
	"errors"
	"fmt"
)

type TaskRepository interface {
	CreateTask(newTask Task) Task
	UpdateTask(updatedTask Task) error
	DeleteTask(taskID int) error
	GetAllTasks() []Task
}

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskRepositoryImpl struct {
	tasks []Task
}

func NewTaskRepository() *TaskRepositoryImpl {
	return &TaskRepositoryImpl{tasks: make([]Task, 0)}
}

func (r *TaskRepositoryImpl) CreateTask(newTask Task) Task {
	newTask.ID = len(r.tasks) + 1
	r.tasks = append(r.tasks, newTask)
	return newTask
}

func (r *TaskRepositoryImpl) UpdateTask(updatedTask Task) error {
	for i, task := range r.tasks {
		fmt.Println("Task ID", task.ID, "Updated Task ID", updatedTask.ID)
		if task.ID == updatedTask.ID {
			r.tasks[i] = updatedTask
			fmt.Println("UPDATED Task", updatedTask)
			return nil
		}
	}
	return errors.New("Task not found")
}

func (r *TaskRepositoryImpl) DeleteTask(taskID int) error {
	for i, task := range r.tasks {
		if task.ID == taskID {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Task not found")
}

func (r *TaskRepositoryImpl) GetAllTasks() []Task {
	return r.tasks
}
