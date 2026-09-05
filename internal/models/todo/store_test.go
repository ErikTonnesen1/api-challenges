package todo_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/ErikTonnesen1/api-challenges/internal/database"
	"github.com/ErikTonnesen1/api-challenges/internal/helpers"
	"github.com/ErikTonnesen1/api-challenges/internal/models/todo"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	db, err := sql.Open("postgres", os.Getenv("APICH_TEST_DSN"))
	if err != nil {
		panic(fmt.Errorf("Couldn't open test DB: %s\n", err.Error()))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Errorf("Couldn't ping db: %s\n", err.Error()))
	}
	testDB = db

	defer testDB.Close()

	os.Exit(m.Run())
}

func cleanTodos(t *testing.T, db *sql.DB) {
	t.Helper()
	t.Cleanup(func() {
		db.Exec("DELETE FROM todos")
	})
}

func TestCRUD(t *testing.T) {
	cleanTodos(t, testDB)

	t.Run("insert ~ should create todo in db", func(t *testing.T) {
		model := database.New(testDB)

		todo := todo.TodoRequest{
			Title: helpers.StringPtr("Write tests"),
			Done:  helpers.BoolPtr(false),
		}

		id, created_at, err := model.Todos.Insert(todo)

		assert.NoError(t, err)
		assert.NotNil(t, id)
		assert.NotNil(t, created_at)

		// insertedItem, err := model.Todos.GetById(id)

		// assert.Equal(t, id, insertedItem.Id)
		// assert.Equal(t, *todo.Title, insertedItem.Title)
		// assert.Equal(t, *todo.Done, insertedItem.Done)
		// assert.NotZero(t, insertedItem.ID)
	})

}
