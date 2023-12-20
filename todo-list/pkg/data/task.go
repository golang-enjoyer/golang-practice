package data

import (
	"errors"
	"fmt"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks []Task

func CreateTask(newTask Task) Task {
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	return newTask
}

func UpdateTask(updatedTask Task) error {
	for i, task := range tasks {
		if task.ID == updatedTask.ID {
			updatedTask.ID = task.ID
			tasks[i] = updatedTask
			return nil
		}
	}
	return errors.New("Task not found")
}

func DeleteTask(taskID int) error {
	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Task not found")
}

func GetAllTasks() []Task {
	return tasks
}
