package todo

import (
	"reflect"
	"testing"

	"github.com/ErikTonnesen1/api-challenges/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	todoService := TodoService{
		2,
		map[int]*TodoItem{
			1: {
				Id:    1,
				Title: "Test Get All",
				Done:  false,
			},
			2: {
				Id:    2,
				Title: "Test Get All Item 2",
				Done:  false,
			},
		},
		nil,
	}

	todoItems := todoService.GetAll(TodoRequest{})

	assert.Equal(t, 2, len(todoItems))
	assert.Equal(t, len(todoService.TodoItems), len(todoItems))
	//Cannot assert on todoItem[key] because go iterates thru maps in random order
	// 		|--> causing intermittent test failures

	// assert.True(t, reflect.DeepEqual(*todoService.TodoItems[1], todoItems[0]))
	// assert.True(t, reflect.DeepEqual(*todoService.TodoItems[2], todoItems[1]))

}

func TestGetAll_WithTitleQuery(t *testing.T) {
	todoService := TodoService{
		id_increment: 2,
		TodoItems: map[int]*TodoItem{
			1: &TodoItem{
				Id:    1,
				Title: "Find me",
				Done:  false,
			},
			2: &TodoItem{
				Id:    1,
				Title: "Don't Find me",
				Done:  false,
			},
		},
	}

	titleQuery := TodoRequest{
		Title: helpers.StringPtr("Find me"),
		Done:  nil,
	}

	foundItems := todoService.GetAll(titleQuery)

	assert.Equal(t, 1, len(foundItems))
	assert.True(t, reflect.DeepEqual(*todoService.TodoItems[1], foundItems[0]))
}

func TestGetAll_WithDoneQuery(t *testing.T) {
	todoService := TodoService{
		id_increment: 2,
		TodoItems: map[int]*TodoItem{
			1: &TodoItem{
				Id:    1,
				Title: "Find me #1",
				Done:  false,
			},
			2: &TodoItem{
				Id:    1,
				Title: "Find me #2",
				Done:  false,
			},
			3: &TodoItem{
				Id:    1,
				Title: "Don't Find me ",
				Done:  true,
			},
		},
	}

	titleQuery := TodoRequest{
		Title: nil,
		Done:  helpers.BoolPtr(false),
	}

	foundItems := todoService.GetAll(titleQuery)

	assert.True(t, len(foundItems) == 2)
	assert.True(t, reflect.DeepEqual(*todoService.TodoItems[1], foundItems[0]))
	assert.True(t, reflect.DeepEqual(*todoService.TodoItems[2], foundItems[1]))
}

func TestGetAll_WithTitleAndDoneQuery(t *testing.T) {
	todoService := TodoService{
		id_increment: 2,
		TodoItems: map[int]*TodoItem{
			1: &TodoItem{
				Id:    1,
				Title: "Find me",
				Done:  false,
			},
			2: &TodoItem{
				Id:    1,
				Title: "Don't Find me",
				Done:  false,
			},
			3: &TodoItem{
				Id:    1,
				Title: "Don't Find me ",
				Done:  true,
			},
		},
	}

	titleQuery := TodoRequest{
		Title: helpers.StringPtr("Find me"),
		Done:  helpers.BoolPtr(false),
	}

	foundItems := todoService.GetAll(titleQuery)

	assert.True(t, len(foundItems) == 1)
	assert.True(t, reflect.DeepEqual(*todoService.TodoItems[1], foundItems[0]))
}

func TestAddItem(t *testing.T) {
	todoService := NewService()
	newTodo := TodoRequest{
		Title: helpers.StringPtr("TestAddItem"),
		Done:  helpers.BoolPtr(false),
	}

	added, err := todoService.AddItem(newTodo)
	assert.True(t, err == nil)
	assert.True(t, added.Id != 0)
	assert.Equal(t, *newTodo.Title, added.Title)
	assert.Equal(t, *newTodo.Done, added.Done)
	assert.Equal(t, 1, len(todoService.TodoItems))

}

func TestGetItem(t *testing.T) {
	todoService := TodoService{
		id_increment: 1,
		TodoItems: map[int]*TodoItem{
			1: {
				Id:    1,
				Title: "Test Get Item",
				Done:  false,
			},
		},
	}

	todoItem, err := todoService.GetItem(1)
	assert.True(t, err == nil)
	assert.True(t, reflect.DeepEqual(*todoService.TodoItems[1], todoItem))
}

func TestGetItem_IdNotFoundThrowsError(t *testing.T) {
	todoService := TodoService{
		id_increment: 1,
		TodoItems: map[int]*TodoItem{
			1: {
				Id:    1,
				Title: "Test Get Item",
				Done:  false,
			},
		},
	}

	todoItem, err := todoService.GetItem(2)
	assert.True(t, err != nil)
	assert.Equal(t, "todo item with id 2 not found", err.Error())
	assert.True(t, reflect.DeepEqual(TodoItem{}, todoItem))
}

func TestToggleDoneService(t *testing.T) {
	todoService := TodoService{
		TodoItems: map[int]*TodoItem{
			1: {
				Id:    1,
				Title: "Test Toggle",
				Done:  false,
			},
		},
	}

	toggled, err := todoService.ToggleDone(1)
	assert.True(t, err == nil)
	assert.Equal(t, true, toggled.Done)
	assert.Equal(t, true, todoService.TodoItems[1].Done)
}

func TestDeleteTodoService(t *testing.T) {
	todoService := TodoService{
		TodoItems: map[int]*TodoItem{
			1: {
				Id:    1,
				Title: "Test Delete",
				Done:  false,
			},
		},
	}

	toBeDeleted := todoService.TodoItems[1]

	deleted, err := todoService.DeleteTodo(1)
	assert.True(t, err == nil)
	assert.True(t, reflect.DeepEqual(*toBeDeleted, deleted))

	_, ok := todoService.TodoItems[1]
	assert.True(t, ok != true)

}

func TestDeleteTodoService_ThrowsErrorIfInvalidId(t *testing.T) {
	todoService := TodoService{
		TodoItems: map[int]*TodoItem{
			1: {
				Id:    1,
				Title: "Test Delete",
				Done:  false,
			},
		},
	}

	_, err := todoService.DeleteTodo(0)
	assert.Equal(t, err.Error(), "No TodoItem found for ID: 0")

}

func TestReplaceTodo(t *testing.T) {
	todoService := TodoService{
		TodoItems: map[int]*TodoItem{
			1: {
				Id:    1,
				Title: "Test Delete",
				Done:  false,
			},
		},
	}

	replacementTodo := TodoRequest{
		Title: helpers.StringPtr("New Title"),
		Done:  helpers.BoolPtr(true),
	}

	new, err := todoService.ReplaceTodo(1, replacementTodo)
	assert.True(t, err == nil)
	assert.Equal(t, *replacementTodo.Title, new.Title)
	assert.Equal(t, *replacementTodo.Done, new.Done)

	todoItemInMap := todoService.TodoItems[1]
	assert.Equal(t, *replacementTodo.Title, todoItemInMap.Title)
	assert.Equal(t, *replacementTodo.Done, todoItemInMap.Done)
}

func TestReplaceTodo_throwsErrorIfIdNotFoundInMap(t *testing.T) {
	todoService := TodoService{
		TodoItems: map[int]*TodoItem{
			1: {
				Id:    1,
				Title: "Test Delete",
				Done:  false,
			},
		},
	}

	replacementTodo := TodoRequest{
		Title: helpers.StringPtr("New Title"),
		Done:  helpers.BoolPtr(true),
	}

	_, err := todoService.ReplaceTodo(0, replacementTodo)
	assert.False(t, err == nil)
	assert.Equal(t, err.Error(), "Could not find item to be replaced with id: 0")
}
