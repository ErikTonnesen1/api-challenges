package todo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RequestValidation() gin.HandlerFunc {
	return func(c *gin.Context) {

		if id := c.Param("id"); id != "" {
			if _, err := strconv.Atoi(id); err != nil {
				throwBadRequestError(c, "ID must be of type int")
				return
			}
		}

		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut {

			var requestTodo TodoItemRequest
			if err := c.ShouldBindJSON(&requestTodo); err != nil {
				throwBadRequestError(c, fmt.Sprintf("Could not parse JSON: %s", err))
				return
			}

			c.Set(TodoRequestContextKey, requestTodo)
		}

		c.Next()
	}
}

func QueryFilterValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		//If passed query params for a GET endpoint
		if c.Query("title") != "" ||
			c.Query("done") != "" {
			doneQueryParam, err := strconv.ParseBool(c.Query("done"))
			if err != nil {
				throwBadRequestError(c, "'done' query param must be of type bool")
				return
			}
			c.Set(TodoRequestQueryFilter, TodoItemRequest{c.Query("title"), doneQueryParam})
		}

		c.Next()
	}
}

func throwBadRequestError(c *gin.Context, errorMsg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": errorMsg,
	})
}
