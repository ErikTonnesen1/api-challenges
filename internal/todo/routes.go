package todo

import (
	"github.com/ErikTonnesen1/api-challenges/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupTodoRoutes(eng *gin.Engine) {
	todoService := NewTodoService()
	todoHandler := NewTodoHandler(todoService)
	todos := eng.Group(TodoUri)
	todos.Use(middleware.Logging(), RequestValidation())

	todos.GET("", QueryFilterValidation(), todoHandler.GetTodos)
	todos.GET("/:id", todoHandler.TodosById)
	todos.POST("", todoHandler.CreateTodo)
	todos.PATCH("/:id", todoHandler.ToggleDone)
	todos.PUT("/:id", todoHandler.ReplaceTodo)
	todos.DELETE("/:id", todoHandler.DeleteTodo)
}
