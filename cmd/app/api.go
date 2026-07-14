package main

import (
	"log"
	"net/http"
	"time"

	todo "github.com/ErikTonnesen1/api-challenges/internal/todo"
	"github.com/gin-gonic/gin"
)

type app struct {
	port string
}

func (a *app) getRoutes() *gin.Engine {
	todoService := todo.NewTodoService()
	todoHandler := todo.NewTodoHandler(todoService)

	routes := gin.Default()

	routes.GET("/todos", todoHandler.GetTodos)
	routes.GET("/todos/:id", todoHandler.TodosById)
	routes.POST("/todos", todoHandler.CreateTodo)
	routes.PATCH("/todos/:id", todoHandler.ToggleDone)
	routes.PUT("/todos/:id", todoHandler.ReplaceTodo)
	routes.DELETE("/todos/:id", todoHandler.DeleteTodo)

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
