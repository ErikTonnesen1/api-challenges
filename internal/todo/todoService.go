package todo

import (
	"fmt"
)

type TodoService struct {
	id_increment int
	TodoItems    []TodoItem
}

func NewTodoService() *TodoService {
	return &TodoService{
		id_increment: 1,
		TodoItems:    []TodoItem{},
	}
}

func (s *TodoService) GetAll() []TodoItem {
	return s.TodoItems
}

func (s *TodoService) AddItem(i TodoItem) (TodoItem, error) {
	if i.Id != 0 {
		return TodoItem{}, fmt.Errorf("Setting an ID is not allowed")
	}
	i.Id = s.id_increment
	s.id_increment++
	s.TodoItems = append(s.TodoItems, i)
	return i, nil
}

func (s *TodoService) GetItem(id int) (TodoItem, error) {
	for _, item := range s.TodoItems {
		if item.Id == id {
			return item, nil
		}
	}
	return TodoItem{}, fmt.Errorf("todo item with id %d not found", id)
}
