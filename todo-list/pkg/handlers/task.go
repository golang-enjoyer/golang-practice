package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks []Task

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(taskIDStr)
	fmt.Println(taskID, err)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	found := false
	for i, task := range tasks {
		if task.ID == taskID {
			found = true

			if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			updatedTask.ID = taskID
			tasks[i] = updatedTask

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}

	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Path[len("/tasks/"):]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	found := false
	for i, task := range tasks {
		if task.ID == taskID {
			found = true

			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Task deleted successfully")
			break
		}
	}

	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
