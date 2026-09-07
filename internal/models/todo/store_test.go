package todo_test

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ErikTonnesen1/api-challenges/internal/database"
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
	model := database.New(testDB)

	t.Run("insert ~ should create todo in db", func(t *testing.T) {
		cleanTodos(t, testDB)
		todo := todo.NewRequest("Write tests", false)

		id, created_at, err := model.Todos.Insert(todo)

		assert.NoError(t, err)
		assert.NotNil(t, id)
		assert.NotNil(t, created_at)

		insertedItem, err := model.Todos.GetById(id)

		assert.Equal(t, id, insertedItem.Id)
		assert.Equal(t, *todo.Title, insertedItem.Title)
		assert.Equal(t, *todo.Done, insertedItem.Done)
		assert.NotZero(t, insertedItem.Id)
	})

	t.Run("Get ~ get by Id should return correct todo", func(t *testing.T) {
		cleanTodos(t, testDB)
		testTodo := insertTestTodo(model)

		getResult, err := model.Todos.GetById(testTodo.Id)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, testTodo.Title, getResult.Title)
		assert.Equal(t, testTodo.Done, getResult.Done)
		assert.NotNil(t, getResult.Id)
		assert.NotNil(t, getResult.CreatedAt)
	})

	t.Run("get ~ get all should return all records in todos", func(t *testing.T) {
		cleanTodos(t, testDB)

		todos := []todo.TodoRequest{
			todo.NewRequest("first task", false),
			todo.NewRequest("second task", false),
			todo.NewRequest("third task", false),
		}

		insertManyTestTodos(todos, model)

		page, err := model.Todos.GetAll()
		if err != nil {
			panic(err)
		}

		assert.Equal(t, 3, len(page))

		for i, v := range todos {
			assert.True(t, assertReqToCreated(v, page[i]))
		}
	})

	t.Run("get ~ get all with empty db should return `no record found ` err", func(t *testing.T) {
		cleanTodos(t, testDB)

		_, err := model.Todos.GetAll()

		assert.Error(t, err)
		assert.ErrorIs(t, todo.ErrRecordNotFound, err)
	})

	// t.Run("get ~ get all with pagination should return that amount of todos", func(t *testing.T) {
	// 	insertManyTestTodos([]todo.TodoRequest{
	// 		todo.NewRequest("first task", false),
	// 		todo.NewRequest("second task", false),
	// 		todo.NewRequest("third task", false),
	// 	}, model)
	//
	// 	cursor, err := model.Todos.GetAllPaginated()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	assert.Equal(t, 3, len(cursor))
	//
	// })

	t.Run("delete ~ delete by Id should delete existing documnets", func(t *testing.T) {
		cleanTodos(t, testDB)
		testTodo := insertTestTodo(model)

		err := model.Todos.Delete(testTodo.Id)
		if err != nil {
			panic(err)
		}
		_, err = model.Todos.GetById(testTodo.Id)
		assert.True(t, errors.Is(err, todo.ErrRecordNotFound))
	})
	t.Run("update ~ updating a document should change existing document", func(t *testing.T) {
		cleanTodos(t, testDB)
		testTodo := insertTestTodo(model)

		update := todo.NewRequest("Finish Tests", true)

		updated, err := model.Todos.Update(testTodo.Id, update)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, testTodo.Id, updated.Id)
		assert.Equal(t, *update.Title, updated.Title)
		assert.Equal(t, *update.Done, updated.Done)
		assert.NotNil(t, updated.CreatedAt)
	})

}

func insertTestTodo(model database.Models) todo.TodoItem {
	testTodo := todo.NewRequest("Write Tests", false)
	id, created_at, err := model.Todos.Insert(testTodo)
	if err != nil {
		panic(err.Error())
	}
	return todo.TodoItem{Id: id, CreatedAt: created_at, Title: *testTodo.Title, Done: *testTodo.Done}
}

func insertManyTestTodos(arr []todo.TodoRequest, model database.Models) {
	for _, v := range arr {
		_, _, err := model.Todos.Insert(v)
		if err != nil {
			panic(err)
		}
	}
}

func assertReqToCreated(req todo.TodoRequest, item todo.TodoItem) bool {
	return *req.Title == item.Title && *req.Done == item.Done
}
