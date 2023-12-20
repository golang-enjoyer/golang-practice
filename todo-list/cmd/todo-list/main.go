package main

import (
	"net/http"
	"todo-list/pkg/handlers"
	"todo-list/pkg/middleware"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Use(middleware.LoggerMiddleware)

	router.HandleFunc("/tasks", handlers.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks", handlers.GetAllTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods(http.MethodDelete)

	http.ListenAndServe("localhost:8080", router)
}
