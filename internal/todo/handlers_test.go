package todo

//
// import (
// 	"bytes"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// type MockTodoService struct {
// 	GetAllResult     []TodoItem
// 	GetItemResult    TodoItem
// 	AddItemResult    TodoItem
// 	ToggleDoneResult TodoItem
// 	DeleteTodoResult TodoItem
// 	ErrorResult      error
// }
//
// func (m *MockTodoService) GetAll() []TodoItem {
// 	return m.GetAllResult
// }
// func (m *MockTodoService) GetItem(id int) (TodoItem, error) {
// 	return m.GetItemResult, m.ErrorResult
// }
//
// func (m *MockTodoService) AddItem(i TodoItem) (TodoItem, error) {
// 	return m.AddItemResult, m.ErrorResult
// }
//
// func (m *MockTodoService) ToggleDone(id int) (TodoItem, error) {
//
// 	m.ToggleDoneResult.Done = !m.ToggleDoneResult.Done
//
// 	return m.ToggleDoneResult, m.ErrorResult
// }
// func (m *MockTodoService) DeleteTodo(id int) (TodoItem, error) {
// 	return m.DeleteTodoResult, m.ErrorResult
// }
//
// func TestGetTodos(t *testing.T) {
// 	//Given
// 	mockTService := MockTodoService{
// 		GetAllResult: []TodoItem{
// 			{
// 				Id:    1,
// 				Title: "Test 1",
// 				Done:  false,
// 			},
// 		},
// 	}
//
// 	h := NewTodoHandler(&mockTService)
//
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest(http.MethodGet, "/todos", nil)
//
// 	//When
// 	h.GetTodos(w, r)
//
// 	//Then
// 	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
// 	assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
//
// 	var returnedTodos []TodoItem
// 	err := json.NewDecoder(w.Body).Decode(&returnedTodos)
// 	if err != nil {
// 		log.Fatal("Returned response from TodoHandler GetAll cannot be parsed into json.")
// 	}
// 	assert.True(t, reflect.DeepEqual(returnedTodos, mockTService.GetAllResult))
// }
//
// func TestGetTodos_Empty(t *testing.T) {
// 	//Given
// 	mockTService := MockTodoService{
// 		GetAllResult: []TodoItem{},
// 	}
//
// 	h := NewTodoHandler(&mockTService)
//
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest(http.MethodGet, "/todos", nil)
//
// 	//When
// 	h.GetTodos(w, r)
//
// 	//Then
// 	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
// 	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
//
// 	var returnedTodos []TodoItem
// 	err := json.NewDecoder(w.Body).Decode(&returnedTodos)
// 	if err != nil {
// 		log.Fatal("Returned response from TodoHandler GetAll cannot be parsed into json.")
// 	}
// 	assert.True(t, reflect.DeepEqual(mockTService.GetAllResult, returnedTodos))
// }
//
// func TestCreateTodo(t *testing.T) {
// 	//Given
// 	mockService := MockTodoService{
// 		AddItemResult: TodoItem{
// 			Id:    1,
// 			Title: "Test Create",
// 			Done:  false,
// 		},
// 	}
//
// 	h := NewTodoHandler(&mockService)
//
// 	reqBody, err := json.Marshal(mockService.AddItemResult)
// 	if err != nil {
// 		log.Printf("Could not encode TodoItem: %v", err)
// 	}
//
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(reqBody))
// 	r.Header.Set("Content-Type", "application/json")
//
// 	h.CreateTodo(w, r)
//
// 	var resultTodo TodoItem
// 	err = json.NewDecoder(w.Body).Decode(&resultTodo)
// 	if err != nil {
// 		log.Printf("Could not decode returned Json: %v", err)
// 	}
//
// 	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
// 	assert.Equal(t, mockService.AddItemResult, resultTodo)
// }
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
