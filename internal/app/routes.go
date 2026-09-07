package app

import (
	"fmt"

	"github.com/ErikTonnesen1/api-challenges/internal/middleware"
	"github.com/ErikTonnesen1/api-challenges/internal/models/todo"
	"github.com/gin-gonic/gin"
)

func (a *Application) StartGin() *gin.Engine {
	engine := gin.Default()
	a.setRoutes(engine)
	return engine
}

func (a *Application) Serve(engine *gin.Engine) error {
	a.Logger.Info(fmt.Sprintf("Starting Gin server on port %s", a.Config.Port))
	return engine.Run(a.Config.Port)
}

func (a *Application) setRoutes(eng *gin.Engine) {
	todoService := todo.NewService(a.Models.Todos)
	todoHandler := todo.NewHandler(todoService)
	todos := eng.Group(todo.TodoUri)
	todos.Use(middleware.Logging(), middleware.RequestValidation())

	todos.GET("", middleware.QueryFilterValidation(), todoHandler.GetTodos)
	todos.GET("/:id", todoHandler.TodosById)
	todos.POST("", todoHandler.CreateTodo)
	todos.PATCH("/:id", todoHandler.ToggleDone)
	todos.PUT("/:id", todoHandler.ReplaceTodo)
	todos.DELETE("/:id", todoHandler.DeleteTodo)
}
