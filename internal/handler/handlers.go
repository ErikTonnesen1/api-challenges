package handler

import (
	"fmt"
	"io"
	"net/http"
)

func Todos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		io.WriteString(w, "Hello from Get TODO's\n")

	case http.MethodPost:
		io.WriteString(w, "Hello from POST TODO's\n")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func TodosById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		io.WriteString(w, fmt.Sprintf("Hello from Get TODO's by ID for: %s\n", r.PathValue("id")))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
