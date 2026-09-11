package todo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ErikTonnesen1/api-challenges/internal/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Set Gin in TestMode for better test log output
func SetupGinInTestMode(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
}

func TestGetTodos(t *testing.T) {
	SetupGinInTestMode(t)
	rw, context := setupRouter(httptest.NewRequest(http.MethodGet, TodoUri, nil))

	//Given
	mockTService := MockTodoService{
		GetAllResult: []TodoItem{
			{
				Id:    1,
				Title: "Test 1",
				Done:  false,
			},
			{
				Id:    2,
				Title: "Test 2",
				Done:  true,
			},
		},
		ErrorResult: nil,
	}

	h := NewHandler(&mockTService)
	//When
	h.GetTodos(context)

	//Then
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))

	var returnedTodos []TodoItem
	err := json.NewDecoder(rw.Body).Decode(&returnedTodos)
	// Useful when you are handling []bytes -- when we already have an io.Reader like rw.Body, it's excessive to use Marshal/Unmarshal
	// err := json.Unmarshal(rw.Body.Bytes(), &returnedTodos)
	if err != nil {
		log.Fatal("Returned response from TodoHandler GetAll cannot be parsed into json.")
	}
	assert.True(t, reflect.DeepEqual(mockTService.GetAllResult, returnedTodos))
}

func TestGetTodos_Empty(t *testing.T) {
	SetupGinInTestMode(t)
	//Given
	rw, context := setupRouter(httptest.NewRequest(http.MethodGet, TodoUri, nil))

	mockTService := MockTodoService{
		GetAllResult: []TodoItem{},
	}

	h := NewHandler(&mockTService)

	//When
	h.GetTodos(context)

	//Then
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))

	var returnedTodos []TodoItem
	err := json.NewDecoder(rw.Body).Decode(&returnedTodos)
	if err != nil {
		log.Fatal("Returned response from TodoHandler GetAll cannot be parsed into json.")
	}
	assert.True(t, reflect.DeepEqual(mockTService.GetAllResult, returnedTodos))
}

func TestCreateTodo(t *testing.T) {
	SetupGinInTestMode(t)
	//Given
	mockService := MockTodoService{
		AddItemResult: TodoItem{
			Id:    1,
			Title: "Test Create",
			Done:  false,
		},
	}

	createItemRequest := TodoRequest{
		Title: helpers.StringPtr("Test Create"),
		Done:  helpers.BoolPtr(false),
	}

	reqBody, err := json.Marshal(createItemRequest)
	if err != nil {
		log.Printf("Could not encode TodoItem: %v", err)
	}
	rw, context := setupRouter(httptest.NewRequest(http.MethodPost, TodoUri, bytes.NewBuffer(reqBody)))
	context.Set(TodoRequestContextKey, createItemRequest)

	h := NewHandler(&mockService)

	h.CreateTodo(context)

	var result TodoItem
	err = json.NewDecoder(rw.Body).Decode(&result)
	if err != nil {
		log.Printf("Could not decode returned Json: %v", err)
	}

	assert.Equal(t, http.StatusCreated, rw.Result().StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))
	assert.ObjectsAreEqual(mockService.AddItemResult, result)
}

func TestTodosById(t *testing.T) {
	SetupGinInTestMode(t)
	mockService := MockTodoService{
		GetItemResult: TodoItem{
			Id:    1,
			Title: "Test Get By ID",
			Done:  false,
		},
	}

	h := NewHandler(&mockService)

	rw, context := setupRouter(httptest.NewRequest(http.MethodGet, "/todos/1", nil))

	h.TodosById(context)

	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))

	var returnedTodo TodoItem
	err := json.NewDecoder(rw.Body).Decode(&returnedTodo)
	if err != nil {
		log.Printf("Error decoding returned Json: %v", err)
	}

	assert.True(t, reflect.DeepEqual(mockService.GetItemResult, returnedTodo))
}

func TestToggleDone(t *testing.T) {
	SetupGinInTestMode(t)
	mockService := MockTodoService{
		ToggleDoneResult: TodoItem{
			Id:    1,
			Title: "Toggle Done",
			Done:  false,
		},
	}

	h := NewHandler(&mockService)

	rw, context := setupRouter(httptest.NewRequest(http.MethodPatch, "/todos/1", nil))

	h.ToggleDone(context)

	var returnedTodo TodoItem
	err := json.NewDecoder(rw.Body).Decode(&returnedTodo)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
	}

	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))
	assert.Equal(t, true, returnedTodo.Done)

}

func TestDeleteTodo(t *testing.T) {
	SetupGinInTestMode(t)
	mockService := MockTodoService{}

	h := NewHandler(&mockService)

	rw, context := setupRouter(httptest.NewRequest(http.MethodDelete, "/todos/1", nil))

	h.DeleteTodo(context)

	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
}

func setupRouter(req *http.Request) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return w, c
}

type MockTodoService struct {
	GetAllResult     []TodoItem
	GetItemResult    TodoItem
	AddItemResult    TodoItem
	ToggleDoneResult TodoItem
	ErrorResult      error
	ReplaceResult    TodoItem
}

func (m *MockTodoService) GetAll(queryFilter TodoRequest) ([]TodoItem, error) {
	return m.GetAllResult, m.ErrorResult
}
func (m *MockTodoService) GetItem(id int) (*TodoItem, error) {
	return &m.GetItemResult, m.ErrorResult
}

func (m *MockTodoService) AddItem(i TodoRequest) (*TodoItem, error) {
	return &m.AddItemResult, m.ErrorResult
}

func (m *MockTodoService) ToggleDone(id int) (*TodoItem, error) {

	m.ToggleDoneResult.Done = !m.ToggleDoneResult.Done

	return &m.ToggleDoneResult, m.ErrorResult
}
func (m *MockTodoService) DeleteTodo(id int) error {
	return m.ErrorResult
}
func (m *MockTodoService) ReplaceTodo(id int, replacement TodoRequest) (*TodoItem, error) {
	return &m.ReplaceResult, m.ErrorResult
}
