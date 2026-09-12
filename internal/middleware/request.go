package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ErikTonnesen1/api-challenges/internal/helpers"
	"github.com/ErikTonnesen1/api-challenges/internal/models/todo"
	"github.com/ErikTonnesen1/api-challenges/internal/validator"
	"github.com/gin-gonic/gin"
)

func RequestValidation() gin.HandlerFunc {
	v := validator.New()

	return func(c *gin.Context) {
		if id := c.Param("id"); id != "" {
			_, err := strconv.Atoi(id)
			v.Check(err == nil, "id", "id must be of type int")
		}

		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			var requestTodo todo.TodoRequest
			err := helpers.ReadJson(c.Writer, c.Request, &requestTodo)
			v.Check(err == nil, "json", fmt.Sprintf("could not parse json: %s", err))
			c.Set(todo.TodoRequestContextKey, requestTodo)
		}
		if !v.Valid() {
			throwBadRequestError(c, v.Errors)
			return
		}
		v.Clear()
		c.Next()
	}
}

func QueryFilterValidation() gin.HandlerFunc {
	v := validator.New()
	return func(c *gin.Context) {
		var titleFilter *string = nil
		var doneFilter *bool = nil

		if tFilter := c.Query("title"); tFilter != "" {
			titleFilter = &tFilter
		}

		if dFilter := c.Query("done"); dFilter != "" {
			doneQueryParam, err := strconv.ParseBool(dFilter)
			v.Check(err == nil, "done", "done must be of type bool")
			if !v.Valid() {
				throwBadRequestError(c, v.Errors)
				return
			}
			doneFilter = &doneQueryParam
		}

		c.Set(todo.TodoRequestQueryFilter, todo.TodoRequest{Title: titleFilter, Done: doneFilter})
		c.Next()
	}
}

func throwBadRequestError(c *gin.Context, errors map[string]string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, errors)
}
