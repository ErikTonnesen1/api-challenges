package todo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var todoById string = "/todos/%s"
var todoByIdUri string = "/todos/:id"

func Test_ifValidId_ThenContinueMiddlewareChain(t *testing.T) {
	w := httptest.NewRecorder()

	nextMwChainCalled := false
	router := getHappyPathMwChain(todoByIdUri, &nextMwChainCalled, RequestValidation())
	validId := "1"
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf(todoById, validId), nil)
	router.ServeHTTP(w, request)

	assert.True(t, nextMwChainCalled)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func Test_ifInvalidId_thenReturn400(t *testing.T) {
	r := httptest.NewRecorder()
	nextMwChainCalled := false

	rtr := getHappyPathMwChain(todoByIdUri, &nextMwChainCalled, RequestValidation())

	incorrectIdType := "id1"
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf(todoById, incorrectIdType), nil)

	rtr.ServeHTTP(r, request)

	assert.False(t, nextMwChainCalled)
	assert.Equal(t, http.StatusBadRequest, r.Code)

}

func getHappyPathMwChain(uri string, happyPathBoolean *bool, handlers ...gin.HandlerFunc) *gin.Engine {
	rtr := gin.New()

	routerHandlers := append(handlers, func(c *gin.Context) { // Could produce side effect of modifying underlying array if capacity allows.
		*happyPathBoolean = true
		c.Status(http.StatusAccepted)
	})

	rtr.GET(
		uri,
		routerHandlers...,
	)

	return rtr
}
