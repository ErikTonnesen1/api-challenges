package todoService

import (
	"fmt"

	"github.com/ErikTonnesen1/api-challenges/internal/dto"
)

type TodoService struct {
	id_increment int
	TodoItems    []dto.TodoItem
}

func NewTodoService() *TodoService {
	return &TodoService{
		id_increment: 1,
		TodoItems:    []dto.TodoItem{},
	}
}

func (s *TodoService) GetAll() []dto.TodoItem {
	return s.TodoItems
}

func (s *TodoService) AddItem(i dto.TodoItem) (dto.TodoItem, error) {
	if i.Id != 0 {
		return dto.TodoItem{}, fmt.Errorf("Setting an ID is not allowed")
	}
	i.Id = s.id_increment
	s.id_increment++
	s.TodoItems = append(s.TodoItems, i)
	return i, nil
}

func (s *TodoService) GetItem(id int) (dto.TodoItem, error) {
	for _, item := range s.TodoItems {
		if item.Id == id {
			return item, nil
		}
	}
	return dto.TodoItem{}, fmt.Errorf("todo item with id %d not found", id)
}
