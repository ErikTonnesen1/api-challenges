package todo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAll(t *testing.T) {

	newTodoItem := TodoItem{
		Title: "TestItem1",
		Done:  false,
	}

	srv := NewTodoService()
	srv.AddItem(newTodoItem)
	handler := NewTodoHandler(srv)

	todoItems := handler.service.GetAll()

	assert.EqualValues(t, len(todoItems), 1)
	assert.Equal(t, newTodoItem.Title, todoItems[0].Title)
	assert.Equal(t, newTodoItem.Done, todoItems[0].Done)
}
