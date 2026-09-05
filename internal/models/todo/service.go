package todo

import (
	"database/sql"
	"fmt"
)

type TodoService struct {
	id_increment int
	TodoItems    map[int]*TodoItem
	Db           *sql.DB
}

func NewService() *TodoService {
	return &TodoService{
		id_increment: 1,
		TodoItems:    make(map[int]*TodoItem),
	}
}

func (s *TodoService) setDB(db *sql.DB) {
	s.Db = db
}

// Side effect: Go iterates over a map in random order
// So each time this method is called, a different order of values will be returned
// To return an ordered set, need to collect sorted key list, then return based on that list
func (s *TodoService) GetAll(queryFilter TodoRequest) []TodoItem {
	todoItems := make([]TodoItem, 0, len(s.TodoItems))

	titleFilterExists := queryFilter.Title != nil
	doneFilterExists := queryFilter.Done != nil

	for _, todoPtr := range s.TodoItems {
		if titleFilterExists && todoPtr.Title != *queryFilter.Title {
			continue
		}
		if doneFilterExists && todoPtr.Done != *queryFilter.Done {
			continue
		}

		todoItems = append(todoItems, *todoPtr)
	}
	return todoItems
}

func (s *TodoService) AddItem(req TodoRequest) (TodoItem, error) {
	newTodo := TodoItem{
		Id:    s.id_increment,
		Title: *req.Title,
		Done:  *req.Done,
	}
	s.id_increment++
	s.TodoItems[newTodo.Id] = &newTodo
	return newTodo, nil
}

func (s *TodoService) GetItem(id int) (TodoItem, error) {
	todoPtr, ok := s.TodoItems[id]
	if !ok {
		return TodoItem{}, fmt.Errorf("todo item with id %d not found", id)
	}
	return *todoPtr, nil
}

func (s *TodoService) ToggleDone(id int) (TodoItem, error) {
	todoPtr, ok := s.TodoItems[id]
	if !ok {
		return TodoItem{}, fmt.Errorf("Invalid ID: %d", id)
	}

	todoPtr.Done = !todoPtr.Done
	return *todoPtr, nil
}

func (s *TodoService) DeleteTodo(id int) (TodoItem, error) {
	deleteItemPtr, ok := s.TodoItems[id]
	if !ok {
		return TodoItem{}, fmt.Errorf("No TodoItem found for ID: %d", id)
	}
	delete(s.TodoItems, id)
	return *deleteItemPtr, nil
}

func (s *TodoService) ReplaceTodo(id int, replacement TodoRequest) (TodoItem, error) {
	todoPtr, ok := s.TodoItems[id]
	if !ok {
		return TodoItem{}, fmt.Errorf("Could not find item to be replaced with id: %d", id)
	}
	*todoPtr = TodoItem{
		Id:    todoPtr.Id,
		Title: *replacement.Title,
		Done:  *replacement.Done,
	}

	return *todoPtr, nil
}
