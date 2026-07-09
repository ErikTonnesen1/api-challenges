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
	ToggleDone(id int) (TodoItem, error)
	DeleteTodo(id int) (TodoItem, error)
}

type todoHandler struct {
	service ITodoService
}

func NewTodoHandler(s ITodoService) *todoHandler {
	return &todoHandler{
		service: s,
	}
}

func (h *todoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.service.GetAll())

}

func (h *todoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
}

func (h *todoHandler) TodosById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
}

func (h *todoHandler) ToggleDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathId := r.PathValue("id")
	id, err := strconv.Atoi(pathId)
	if err != nil {
		http.Error(w, fmt.Sprintf("ID could not be parsed: %s", pathId), http.StatusBadRequest)
	}
	toggledItem, err := h.service.ToggleDone(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toggledItem)
}

func (h *todoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathId := r.PathValue("id")
	id, err := strconv.Atoi(pathId)
	if err != nil {
		http.Error(w, fmt.Sprintf("ID could not be parsed: %s", pathId), http.StatusBadRequest)
	}
	deletedItem, err := h.service.DeleteTodo(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deletedItem)
}
