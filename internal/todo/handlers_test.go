package todo

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockTodoService struct {
	mockItems   []TodoItem
	mockedError error
}

func (m *MockTodoService) GetAll() []TodoItem {
	return m.mockItems
}
func (m *MockTodoService) GetItem(id int) (TodoItem, error) {
	return m.mockItems[0], m.mockedError
}

func (m *MockTodoService) AddItem(i TodoItem) (TodoItem, error) {
	return m.mockItems[0], m.mockedError
}

func (m *MockTodoService) ToggleDone(id int) (TodoItem, error) {
	return m.mockItems[0], m.mockedError
}
func (m *MockTodoService) DeleteTodo(id int) (TodoItem, error) {
	return m.mockItems[0], m.mockedError
}

func TestGetTodos(t *testing.T) {
	//Given
	mockTService := MockTodoService{
		[]TodoItem{
			{
				Id:    1,
				Title: "TestItem1",
				Done:  false,
			},
		},
		nil,
	}

	h := NewTodoHandler(&mockTService)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/todos", nil)

	//When
	h.GetTodos(w, r)

	//Then
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json")

	var returnedTodos []TodoItem
	log.Printf("Returned from MockService GetTodos: %+v", h.service.GetAll())

	err := json.NewDecoder(w.Body).Decode(returnedTodos)
	if err != nil {
		log.Fatal("Returned response from TodoHandler GetAll cannot be parsed into json.")
	}

	assert.True(t, reflect.DeepEqual(returnedTodos, mockTService.mockItems))
}
