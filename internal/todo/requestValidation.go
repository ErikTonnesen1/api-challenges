package todo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequestValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut {
			var requestTodo TodoItem
			if err := c.ShouldBindJSON(&requestTodo); err != nil {
				throwBadRequestError("Could not parse JSON", c)
				return
			}
			if requestTodo.Id != 0 {
			}
		}

		c.Next()
	}
}

func throwBadRequestError(errorMsg string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errorMsg,
	})
	return
}
