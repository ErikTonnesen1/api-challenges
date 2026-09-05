package todo

import "database/sql"

type TodoItem struct {
	Id    int    `json:"id"`
	Title string `json:"title" binding:"required"`
	Done  bool   `json:"done"`
}

type TodoModel struct {
	Db *sql.DB
}
