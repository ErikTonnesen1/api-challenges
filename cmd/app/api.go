package main

import (
	"log"

	"github.com/ErikTonnesen1/api-challenges/internal/todo"
	"github.com/gin-gonic/gin"
)

type app struct {
	port string
}

func (a *app) getRoutes() *gin.Engine {
	engine := gin.Default()
	todo.SetupTodoRoutes(engine)
	return engine
}

func (a *app) serveGin(engine *gin.Engine) error {
	log.Printf("Starting Gin server on port %s", a.port)
	return engine.Run(a.port)
}
