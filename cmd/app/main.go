package main

import (
	"log"
)

func main() {
	application := app{
		port:    ":8080",
		ginPort: ":8081",
	}

	engine := application.getRoutes()

	if ginErr := application.serveGin(engine); ginErr != nil {
		log.Fatalf("Error starting Gin server: %v", ginErr)
	}
}
