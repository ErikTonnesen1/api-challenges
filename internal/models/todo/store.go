package todo

import (
	"database/sql"
	"errors"
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

var (
	ErrRecordNotFound = errors.New("record not found")
)

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

func (tm *TodoModel) Delete(id int) error {
	query := `
	DELETE FROM todos
	WHERE id = $1
	`

	res, err := tm.Db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (tm *TodoModel) Update(id int, update TodoRequest) (*TodoItem, error) {
	query := `
		UPDATE todos 
		SET title = $1, done = $2
		WHERE id = $3
		RETURNING id, title, done, created_at
	`

	var updated TodoItem

	err := tm.Db.QueryRow(query, update.Title, update.Done, id).Scan(
		&updated.Id,
		&updated.Title,
		&updated.Done,
		&updated.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (tm *TodoModel) GetById(id int) (*TodoItem, error) {
	query := `
	SELECT id, created_at, title, done
	FROM todos
	WHERE id = $1
	`
	var todo TodoItem

	err := tm.Db.QueryRow(query, id).Scan(
		&todo.Id,
		&todo.CreatedAt,
		&todo.Title,
		&todo.Done)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}
	return &todo, nil
}
