package todo

import ()

type TodoService struct {
	db TodoModel
}

func NewService(model TodoModel) *TodoService {
	return &TodoService{
		db: model,
	}
}

func (s *TodoService) GetAll(queryFilter TodoRequest) ([]TodoItem, error) {
	todos, err := s.db.GetAll(queryFilter)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *TodoService) AddItem(req TodoRequest) (int, error) {
	id, err := s.db.Insert(req)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TodoService) GetItem(id int) (*TodoItem, error) {
	todo, err := s.db.GetById(id)
	if err != nil {
		return nil, err

	}
	return todo, nil
}

func (s *TodoService) ToggleDone(id int) (*TodoItem, error) {
	todo, err := s.db.Toggle(id)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *TodoService) DeleteTodo(id int) error {
	err := s.db.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoService) ReplaceTodo(id int, replacement TodoRequest) (*TodoItem, error) {
	todo, err := s.db.Update(id, replacement)
	if err != nil {
		return nil, err
	}
	return todo, nil
}
