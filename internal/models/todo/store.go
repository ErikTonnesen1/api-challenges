package todo

import (
	"database/sql"
	"time"
)

type TodoItem struct {
	Id        int       `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoModel struct {
	Db *sql.DB
}

func (tm *TodoModel) Insert(t TodoRequest) (id int, created_at time.Time, err error) {
	query := `
		INSERT INTO todos (title, done)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err = tm.Db.QueryRow(query, t.Title, t.Done).Scan(&id, &created_at)
	if err != nil {
		return 0, time.Time{}, err
	}
	return id, created_at, nil
}
