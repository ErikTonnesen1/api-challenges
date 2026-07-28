package main

import (
	"github.com/ErikTonnesen1/api-challenges/internal/middleware"
	"github.com/ErikTonnesen1/api-challenges/internal/todo"
	"github.com/gin-gonic/gin"
	"log"
)

type app struct {
	port string
}

func (a *app) getRoutes() *gin.Engine {
	todoService := todo.NewTodoService()
	todoHandler := todo.NewTodoHandler(todoService)

	routes := gin.Default()
	todos := routes.Group(todo.TodoUri)
	todos.Use(middleware.Logging(), todo.RequestValidation())

	todos.GET("", todoHandler.GetTodos)
	todos.GET("/:id", todoHandler.TodosById)
	todos.POST("", todoHandler.CreateTodo)
	todos.PATCH("/:id", todoHandler.ToggleDone)
	todos.PUT("/:id", todoHandler.ReplaceTodo)
	todos.DELETE("/:id", todoHandler.DeleteTodo)

	return routes
}

func (a *app) serveGin(engine *gin.Engine) error {
	log.Printf("Starting Gin server on port %s", a.port)
	return engine.Run(a.port)
}
