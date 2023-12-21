package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"todo-list/pkg/data"
	"todo-list/pkg/handlers"
	"todo-list/pkg/middleware"
)

func main() {
	router := mux.NewRouter()

	router.Use(middleware.LoggerMiddleware)

	taskRepo := data.NewTaskRepository()

	handler := handlers.NewTaskHandler(taskRepo)

	router.HandleFunc("/tasks", handler.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks", handler.GetAllTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", handler.UpdateTask).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", handler.DeleteTask).Methods(http.MethodDelete)

	http.ListenAndServe("localhost:8080", router)
}
