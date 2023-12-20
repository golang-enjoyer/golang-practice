package main

import (
	"net/http"
	"todo-list/pkg/handlers"
)

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateTask(w, r)
		case http.MethodPut:
			handlers.UpdateTask(w, r)
		case http.MethodDelete:
			handlers.DeleteTask(w, r)
		case http.MethodGet:
			handlers.GetAllTasks(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
