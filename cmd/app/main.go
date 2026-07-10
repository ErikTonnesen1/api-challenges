package main

import (
	"log"
	"os"
)

func main() {
	application := app{
		port: ":8080",
	}

	if err := application.serve(application.mount()); err != nil {
		log.Fatalf("Error starting server: %v", err)
		os.Exit(1)
	}
}
