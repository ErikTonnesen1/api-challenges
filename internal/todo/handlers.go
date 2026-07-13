package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *todoHandler) GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetAll())

}

func (h *todoHandler) CreateTodo(c *gin.Context) {
	var newItem TodoItem
	err := json.NewDecoder(c.Request.Body).Decode(&newItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}
	addedItem, err := h.service.AddItem(newItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, addedItem)
}

func (h *todoHandler) TodosById(c *gin.Context) {
	requestedId := c.Param("id")
	id, err := strconv.Atoi(requestedId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Id must be of type int: %s", requestedId),
		})
		return
	}
	todoItem, err := h.service.GetItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, todoItem)
	}
}

func (h *todoHandler) ToggleDone(c *gin.Context) {
	pathId := c.Param("id")
	id, err := strconv.Atoi(pathId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("ID could not be parsed %s", pathId),
		})
		return
	}
	toggledItem, err := h.service.ToggleDone(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, toggledItem)
}

func (h *todoHandler) DeleteTodo(c *gin.Context) {
	pathId := c.Param("id")
	id, err := strconv.Atoi(pathId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("ID could not be parsed: %s", pathId),
		})
		return
	}
	deletedItem, err := h.service.DeleteTodo(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, deletedItem)
}
