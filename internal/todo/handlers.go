package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Use an interface as expected Handler param in order to inject test/mock services
// Apparently, the definition of the Interface usually lives in the consumer, while the
// concrete struct will live in the service package
type ITodoService interface {
	GetAll() []TodoItem
	GetItem(id int) (TodoItem, error)
	AddItem(i TodoItem) (TodoItem, error)
}

type todoHandler struct {
	service ITodoService
}

func NewTodoHandler(s ITodoService) *todoHandler {
	return &todoHandler{
		service: s,
	}
}

func (h *todoHandler) Todos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(h.service.GetAll())

	case http.MethodPost:
		var newItem TodoItem
		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		addedItem, err := h.service.AddItem(newItem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(addedItem)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func (h *todoHandler) TodosById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		requestedId := r.PathValue("id")
		id, err := strconv.Atoi(requestedId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Id must be of type int: %s", requestedId), http.StatusBadRequest)
		}
		todoItem, err := h.service.GetItem(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(todoItem)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
