package todo

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ErikTonnesen1/api-challenges/internal/helpers"
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

func NewRequest(title string, done bool) TodoRequest {
	return TodoRequest{
		Title: helpers.StringPtr(title),
		Done:  helpers.BoolPtr(done),
	}
}

type TodoServicer interface {
	GetAll(TodoRequest) ([]TodoItem, error)
	GetItem(id int) (*TodoItem, error)
	AddItem(i TodoRequest) (*TodoItem, error)
	ToggleDone(id int) (*TodoItem, error)
	DeleteTodo(id int) error
	ReplaceTodo(id int, replacement TodoRequest) (*TodoItem, error)
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

	var (
		todos []TodoItem
		err   error
	)

	if filterExists {
		todos, err = h.service.GetAll(queryFilter.(TodoRequest))
	} else {
		todos, err = h.service.GetAll(TodoRequest{})
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, todos)
}

func (h *todoHandler) CreateTodo(c *gin.Context) {
	newItem := c.MustGet(TodoRequestContextKey).(TodoRequest)
	todo, err := h.service.AddItem(newItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, *todo)
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
		c.JSON(http.StatusOK, *todoItem)
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
	c.JSON(http.StatusOK, *toggledItem)
}

func (h *todoHandler) DeleteTodo(c *gin.Context) {
	pathId := c.Param("id")
	id, err := strconv.Atoi(pathId)
	//TODO: Add getParamId helper to return error if id is invalid
	err = h.service.DeleteTodo(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		default:
			c.Status(http.StatusBadRequest)

		}
		return
	}
	c.Status(http.StatusOK)
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
	c.JSON(http.StatusOK, *replacedItem)

}
