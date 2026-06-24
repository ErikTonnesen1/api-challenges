package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ErikTonnesen1/api-challenges/internal/todo"
)

type api struct {
	port string
}

func (a *api) mount() http.Handler {
	mux := http.NewServeMux()

	service := todo.NewTodoService()
	handler := todo.NewTodoHandler(service)

	mux.HandleFunc("GET /todos", handler.Todos)
	mux.HandleFunc("GET /todos/{id}", handler.TodosById)

	mux.HandleFunc("POST /todos", handler.Todos)

	mux.HandleFunc("PATCH /todos", handler.Todos)

	mux.HandleFunc("DELETE /todos", handler.Todos)

	return mux
}

func (a *api) serve(h http.Handler) error {
	srv := &http.Server{
		Addr:              a.port,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       time.Minute,
	}
	log.Printf("server started on port: %s", a.port)

	return srv.ListenAndServe()
}
