package server

import (
	"fmt"
	"github.com/ErikTonnesen1/api-challenges/internal/handler"
	"log"
	"net/http"
)

func StartServer() {
	fmt.Println("Starting server...")

	http.HandleFunc("/todos", handler.Todos)
	http.HandleFunc("/todos/{id}", handler.TodosById)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
