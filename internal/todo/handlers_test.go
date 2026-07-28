package todo

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var todosByIdUri string = "/todos/:id"

// Set Gin in TestMode for better test log output
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestGetTodos(t *testing.T) {
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
	}

	h := NewTodoHandler(&mockTService)
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
	//Given
	rw, context := setupRouter(httptest.NewRequest(http.MethodGet, TodoUri, nil))

	mockTService := MockTodoService{
		GetAllResult: []TodoItem{},
	}

	h := NewTodoHandler(&mockTService)

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
	//Given
	mockService := MockTodoService{
		AddItemResult: TodoItem{
			Id:    1,
			Title: "Test Create",
			Done:  false,
		},
	}

	createItemRequest := TodoItemRequest{
		Title: "Test Create",
		Done:  false,
	}

	reqBody, err := json.Marshal(createItemRequest)
	if err != nil {
		log.Printf("Could not encode TodoItem: %v", err)
	}
	rw, context := setupRouter(httptest.NewRequest(http.MethodPost, TodoUri, bytes.NewBuffer(reqBody)))
	context.Set(TodoRequestContextKey, createItemRequest)

	h := NewTodoHandler(&mockService)

	h.CreateTodo(context)

	var resultTodo TodoItem
	err = json.NewDecoder(rw.Body).Decode(&resultTodo)
	if err != nil {
		log.Printf("Could not decode returned Json: %v", err)
	}

	assert.Equal(t, http.StatusCreated, rw.Result().StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))
	assert.Equal(t, mockService.AddItemResult, resultTodo)
	assert.Equal(t, createItemRequest.Title, resultTodo.Title)
	assert.Equal(t, createItemRequest.Done, resultTodo.Done)
}

//
// func TestTodosById(t *testing.T) {
// 	mockService := MockTodoService{
// 		GetItemResult: TodoItem{
// 			Id:    1,
// 			Title: "Test Get By ID",
// 			Done:  false,
// 		},
// 	}
//
// 	h := NewTodoHandler(&mockService)
//
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
// 	r.SetPathValue("id", "1")
//
// 	h.TodosById(w, r)
//
// 	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
// 	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
//
// 	var returnedTodo TodoItem
// 	err := json.NewDecoder(w.Body).Decode(&returnedTodo)
// 	if err != nil {
// 		log.Printf("Error decoding returned Json: %v", err)
// 	}
//
// 	assert.True(t, reflect.DeepEqual(mockService.GetItemResult, returnedTodo))
// }
//
// func TestToggleDone(t *testing.T) {
// 	mockService := MockTodoService{
// 		ToggleDoneResult: TodoItem{
// 			Id:    1,
// 			Title: "Toggle Done",
// 			Done:  false,
// 		},
// 	}
//
// 	h := NewTodoHandler(&mockService)
//
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest(http.MethodPatch, "/todos/1", nil)
// 	r.SetPathValue("id", "1")
//
// 	h.ToggleDone(w, r)
//
// 	var returnedTodo TodoItem
// 	err := json.NewDecoder(w.Body).Decode(&returnedTodo)
// 	if err != nil {
// 		log.Printf("Error decoding json: %v", err)
// 	}
//
// 	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
// 	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
// 	assert.Equal(t, true, returnedTodo.Done)
//
// }
//
// func TestDeleteTodo(t *testing.T) {
// 	mockService := MockTodoService{
// 		DeleteTodoResult: TodoItem{
// 			Id:    1,
// 			Title: "Test Delete",
// 			Done:  false,
// 		},
// 	}
//
// 	h := NewTodoHandler(&mockService)
//
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
// 	r.SetPathValue("id", "1")
//
// 	h.DeleteTodo(w, r)
//
// 	var returnedTodo TodoItem
// 	err := json.NewDecoder(w.Body).Decode(&returnedTodo)
// 	if err != nil {
// 		log.Printf("Error decoding json: %v", err)
// 	}
//
// 	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
// 	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
// 	assert.True(t, reflect.DeepEqual(mockService.DeleteTodoResult, returnedTodo))
// }

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
	DeleteTodoResult TodoItem
	ErrorResult      error
	ReplaceResult    TodoItem
}

func (m *MockTodoService) GetAll() []TodoItem {
	return m.GetAllResult
}
func (m *MockTodoService) GetItem(id int) (TodoItem, error) {
	return m.GetItemResult, m.ErrorResult
}

func (m *MockTodoService) AddItem(i TodoItemRequest) (TodoItem, error) {
	return m.AddItemResult, m.ErrorResult
}

func (m *MockTodoService) ToggleDone(id int) (TodoItem, error) {

	m.ToggleDoneResult.Done = !m.ToggleDoneResult.Done

	return m.ToggleDoneResult, m.ErrorResult
}
func (m *MockTodoService) DeleteTodo(id int) (TodoItem, error) {
	return m.DeleteTodoResult, m.ErrorResult
}
func (m *MockTodoService) ReplaceTodo(id int, replacement TodoItemRequest) (TodoItem, error) {
	return m.ReplaceResult, m.ErrorResult
}
