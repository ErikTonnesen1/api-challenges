package todo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_ifInvalidId_thenReturn400(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	incorrectIdType := "id1"
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/todos/%s", incorrectIdType), nil)
	r.Header.Set("Content-Type", "application/json")

	c.Request = r

	mw := RequestValidation()
	mw(c)

	assert.Equal(t, c.Request.Response.StatusCode, http.StatusBadRequest)

}
