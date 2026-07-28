package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var todoById string = "/todos/%s"
var todoByIdUri string = "/todos/:id"

type testState struct {
	nextHandlerCalled    bool
	requestBodyValidated bool
}

func Test_ifValidId_ThenContinueMiddlewareChain(t *testing.T) {
	w := httptest.NewRecorder()

	router, state := getHappyPathMwChain(http.MethodGet, todoByIdUri, RequestValidation())
	validId := "1"
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf(todoById, validId), nil)
	router.ServeHTTP(w, request)

	assert.True(t, state.nextHandlerCalled)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func Test_ifInvalidId_thenReturn400(t *testing.T) {
	r := httptest.NewRecorder()

	rtr, state := getHappyPathMwChain(http.MethodGet, todoByIdUri, RequestValidation())

	incorrectIdType := "id1"
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf(todoById, incorrectIdType), nil)

	rtr.ServeHTTP(r, request)

	assert.False(t, state.nextHandlerCalled)
	assert.Equal(t, http.StatusBadRequest, r.Code)

	var jsonResponse map[string]any
	json.NewDecoder(r.Body).Decode(&jsonResponse)

	assert.Equal(t, "ID must be of type int", jsonResponse["error"])

}

func Test_ifPostMethod_andRequestBodyIsValid_thenReturn200(t *testing.T) {
	//Given
	rw := httptest.NewRecorder()

	validTodoItemRequest := TodoItemRequest{
		Title: "Test Validation Json Handling",
		Done:  false,
	}

	rtr, state := getHappyPathMwChain(http.MethodPost, TodoUri, RequestValidation())

	validRequestAsJson, _ := json.Marshal(validTodoItemRequest)

	request := httptest.NewRequest(http.MethodPost, TodoUri, bytes.NewBuffer(validRequestAsJson))
	rtr.ServeHTTP(rw, request)

	assert.Equal(t, rw.Code, http.StatusAccepted)
	assert.True(t, state.requestBodyValidated)
}

func Test_ifPostMethod_andRequestBodyIsInValid_thenReturn400(t *testing.T) {
	//Given
	rw := httptest.NewRecorder()

	invalidTodoItemRequest := map[string]any{
		"RandomField": "RandomFieldValue",
	}

	rtr, state := getHappyPathMwChain(http.MethodPost, TodoUri, RequestValidation())

	invalidRequestAsJson, _ := json.Marshal(invalidTodoItemRequest)

	request := httptest.NewRequest(http.MethodPost, TodoUri, bytes.NewBuffer(invalidRequestAsJson))
	rtr.ServeHTTP(rw, request)

	assert.Equal(t, rw.Code, http.StatusBadRequest)
	assert.False(t, state.requestBodyValidated)

	var jsonResponse map[string]any
	json.NewDecoder(rw.Body).Decode(&jsonResponse)
	assert.True(t, strings.Contains(jsonResponse["error"].(string), "Could not parse JSON:"))
}

func getHappyPathMwChain(httpMethod string, uri string, handlers ...gin.HandlerFunc) (*gin.Engine, *testState) {
	state := &testState{}
	rtr := gin.New()

	routerHandlers := append(handlers, func(c *gin.Context) {
		state.nextHandlerCalled = true
		c.Status(http.StatusAccepted)
		if httpMethod == http.MethodPost {
			_, ok := c.Get(TodoRequestContextKey)
			state.requestBodyValidated = ok
		}
	})

	switch httpMethod {
	case http.MethodGet:
		rtr.GET(
			uri,
			routerHandlers...,
		)
	case http.MethodPost:
		rtr.POST(
			uri,
			routerHandlers...,
		)
	}

	return rtr, state
}
