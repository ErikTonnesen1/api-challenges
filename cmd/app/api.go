package main

import (
	"net/http"
	"time"

	todo "github.com/ErikTonnesen1/api-challenges/internal/todo"
)

type app struct {
	port string
}

func (a *app) mount() *http.ServeMux {
	todoService := todo.NewTodoService()
	todoHandler := todo.NewTodoHandler(todoService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos/", todoHandler.GetTodos)
	mux.HandleFunc("GET /todos/{id}", todoHandler.TodosById)
	mux.HandleFunc("PATCH /todos/{id}", todoHandler.ToggleDone)
	mux.HandleFunc("POST /todos/", todoHandler.CreateTodo)
	mux.HandleFunc("DELETE /todos/{id}", todoHandler.DeleteTodo)

	return mux
}

func (a *app) serve(multiplexer *http.ServeMux) {
	server := http.Server{
		Addr:         a.port,
		Handler:      multiplexer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	server.ListenAndServe()
}
