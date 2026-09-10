package todo

import (
	"database/sql"
	"errors"
	"fmt"
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

func (tm *TodoModel) Insert(t TodoRequest) (id int, err error) {
	query := `
		INSERT INTO todos (title, done)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err = tm.Db.QueryRow(query, t.Title, t.Done).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (tm *TodoModel) Toggle(id int) (TodoItem, error) {
	query := `
		UPDATE todos
		SET done = NOT done
		WHERE id = $1
		RETURNING id, title, done, created_at
	`
	var t TodoItem
	err := tm.Db.QueryRow(query, id).Scan(&t.Id, &t.Title, &t.Done, &t.CreatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return TodoItem{}, ErrRecordNotFound
		default:
			return TodoItem{}, nil
		}
	}
	return t, nil
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

func (tm *TodoModel) GetAll(r TodoRequest) ([]TodoItem, error) {
	query := `
	SELECT id, created_at, title, done
	FROM todos
	ORDER BY id
	`

	rows, err := tm.Db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []TodoItem

	for rows.Next() {
		var todo TodoItem

		err := rows.Scan(
			&todo.Id,
			&todo.CreatedAt,
			&todo.Title,
			&todo.Done,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if len(todos) == 0 {
		return nil, ErrRecordNotFound
	}

	if queryErr := rows.Err(); queryErr != nil {
		return nil, queryErr
	}

	return todos, nil
}

func (tm TodoModel) buildGetTodoQuery(t TodoRequest) (string, []interface{}) {
	query := "SELECT id, title, done, created_at FROM todos WHERE 1=1"
	args := []interface{}{}
	argPos := 1
	if t.Title != nil {
		query += fmt.Sprintf("AND title = $%d", argPos)
		args = append(args, *t.Title)
		argPos++
	}
	if t.Done != nil {
		query += fmt.Sprintf("AND done = $%d", argPos)
		args = append(args, *t.Done)
		argPos++
	}
	return query, args
}
