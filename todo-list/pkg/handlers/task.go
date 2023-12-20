// handlers/task.go

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"todo-list/pkg/data"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask data.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdTask := data.CreateTask(newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(taskIDStr)
	fmt.Println(taskID, err)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask data.Task
	foundErr := data.UpdateTask(updatedTask)
	if foundErr != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Path[len("/tasks/"):]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = data.DeleteTask(taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Task deleted successfully")
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := data.GetAllTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
