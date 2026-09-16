package database

import (
	"database/sql"

	"github.com/ErikTonnesen1/api-challenges/internal/models/todo"
)

type Models struct {
	Todos todo.TodoModel
}

func New(db *sql.DB) Models {
	return Models{
		Todos: todo.TodoModel{Db: db},
	}
}
