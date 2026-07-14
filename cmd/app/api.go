package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ErikTonnesen1/api-challenges/internal/todo"
	"github.com/gin-gonic/gin"
)

type app struct {
	port string
}

func (a *app) getRoutes() *gin.Engine {
	todoService := todo.NewTodoService()
	todoHandler := todo.NewTodoHandler(todoService)

	routes := gin.Default()
	todos := routes.Group("/todos")
	todos.Use(todo.RequestValidation())

	todos.GET("", todoHandler.GetTodos)
	todos.GET("/:id", todoHandler.TodosById)
	todos.POST("", todoHandler.CreateTodo)
	todos.PATCH("/:id", todoHandler.ToggleDone)
	todos.PUT("/:id", todoHandler.ReplaceTodo)
	todos.DELETE("/:id", todoHandler.DeleteTodo)

	return routes
}

func (a *app) serve(multiplexer *http.ServeMux) error {
	server := http.Server{
		Addr:         a.port,
		Handler:      multiplexer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting server on port %s", a.port)
	return server.ListenAndServe()
}

func (a *app) serveGin(engine *gin.Engine) error {
	log.Printf("Starting Gin server on port %s", a.port)
	return engine.Run(a.port)
}
