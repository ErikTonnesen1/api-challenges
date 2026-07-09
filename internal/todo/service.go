package todo

import (
	"fmt"
)

type TodoService struct {
	id_increment int
	TodoItems    map[int]*TodoItem
}

func NewTodoService() *TodoService {
	return &TodoService{
		id_increment: 1,
		TodoItems:    make(map[int]*TodoItem),
	}
}

// Side effect: Go iterates over a map in random order
// SO each time this method is called, a different order of values will be returned
// To return an ordered set, need to collect sorted key list, then return based on that list
func (s *TodoService) GetAll() []TodoItem {
	todoItems := make([]TodoItem, 0, len(s.TodoItems))

	for _, item := range s.TodoItems {
		todoItems = append(todoItems, *item)
	}
	return todoItems
}

func (s *TodoService) AddItem(i TodoItem) (TodoItem, error) {
	if i.Id != 0 {
		return TodoItem{}, fmt.Errorf("Setting an ID is not allowed")
	}
	i.Id = s.id_increment
	s.id_increment++
	s.TodoItems[i.Id] = &i
	return i, nil
}

func (s *TodoService) GetItem(id int) (TodoItem, error) {
	todoItem, ok := s.TodoItems[id]
	if !ok {
		return TodoItem{}, fmt.Errorf("todo item with id %d not found", id)
	}
	return *todoItem, nil
}

func (s *TodoService) ToggleDone(id int) (TodoItem, error) {
	todoItem, ok := s.TodoItems[id]
	if !ok {
		return TodoItem{}, fmt.Errorf("Invalid ID", id)
	}

	todoItem.Done = !todoItem.Done
	return *todoItem, nil
}

func (s *TodoService) DeleteTodo(id int) (TodoItem, error) {
	deleteItem, ok := s.TodoItems[id]
	if !ok {
		return TodoItem{}, fmt.Errorf("No TodoItem found for ID: %s", id)
	}
	delete(s.TodoItems, id)
	return *deleteItem, nil
}
