package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ErikTonnesen1/api-challenges/internal/handler"
	todoService "github.com/ErikTonnesen1/api-challenges/internal/service"
)

func StartServer() {
	fmt.Println("Starting server...")

	todoService := todoService.NewTodoService()
	todoHandler := handler.NewTodoHandler(todoService)

	http.HandleFunc("/todos", todoHandler.Todos)
	http.HandleFunc("/todos/{id}", todoHandler.TodosById)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
