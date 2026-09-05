package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Use an interface as expected Handler param in order to inject test/mock services
// Apparently, the definition of the Interface usually lives in the consumer, while the
// concrete struct will live in the service package

var TodoUri string = "/todos"
var TodoRequestContextKey string = "todoRequest"
var TodoRequestQueryFilter string = "todoRequestFilter"

// Uses pointers to have explicit optional fields
type TodoRequest struct {
	Title *string `json:"title" binding:"required"`
	Done  *bool   `json:"done"`
}

type TodoServicer interface {
	GetAll(TodoRequest) []TodoItem
	GetItem(id int) (TodoItem, error)
	AddItem(i TodoRequest) (TodoItem, error)
	ToggleDone(id int) (TodoItem, error)
	DeleteTodo(id int) (TodoItem, error)
	ReplaceTodo(id int, replacement TodoRequest) (TodoItem, error)
}

type todoHandler struct {
	service TodoServicer
}

func NewHandler(s TodoServicer) *todoHandler {
	return &todoHandler{
		service: s,
	}
}

func (h *todoHandler) GetTodos(c *gin.Context) {
	queryFilter, filterExists := c.Get(TodoRequestQueryFilter)
	if filterExists {
		c.JSON(http.StatusOK, h.service.GetAll(queryFilter.(TodoRequest)))
	} else {
		c.JSON(http.StatusOK, h.service.GetAll(TodoRequest{}))
	}
}

func (h *todoHandler) CreateTodo(c *gin.Context) {
	newItem := c.MustGet(TodoRequestContextKey).(TodoRequest)
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
	id, _ := strconv.Atoi(requestedId)
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
	id, _ := strconv.Atoi(pathId)
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
	deletedItem, err := h.service.DeleteTodo(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, deletedItem)
}

func (h *todoHandler) ReplaceTodo(c *gin.Context) {
	pathId := c.Param("id")
	id, _ := strconv.Atoi(pathId)
	replacement := c.MustGet(TodoRequestContextKey).(TodoRequest)
	replacedItem, err := h.service.ReplaceTodo(id, replacement)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, replacedItem)

}
