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
	mockItems   map[int]TodoItem
	mockedError error
}

func (m *MockTodoService) GetAll() []TodoItem {
	todoSlice := make([]TodoItem, 0, len(m.mockItems))
	for _, item := range m.mockItems {
		todoSlice = append(todoSlice, item)
	}
	return todoSlice
}
func (m *MockTodoService) GetItem(id int) (TodoItem, error) {
	item, ok := m.mockItems[id]
	if !ok {
		return TodoItem{}, m.mockedError
	}
	return item, nil
}

func (m *MockTodoService) AddItem(i TodoItem) (TodoItem, error) {
	return i, m.mockedError
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
		map[int]TodoItem{
			1: {
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
	err := json.NewDecoder(w.Body).Decode(&returnedTodos)
	if err != nil {
		log.Fatal("Returned response from TodoHandler GetAll cannot be parsed into json.")
	}
	assert.True(t, reflect.DeepEqual(returnedTodos[0], mockTService.mockItems[1]))
}

func TestCreateTodo(t *testing.T) {
	//Given
	mockService := MockTodoService{
		map[int]TodoItem{
			1: {
				Id:    1,
				Title: "Test Item 1",
				Done:  false,
			},
		},
		nil,
	}

	h := NewTodoHandler(&mockService)

	newItem := TodoItem{}

	//When
	h.service.AddItem()

}
