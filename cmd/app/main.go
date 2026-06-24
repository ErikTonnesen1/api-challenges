package main

import (
	"log"
	"os"
)

func main() {
	todoApplication := api{
		port: "8080",
	}

	if err := todoApplication.serve(todoApplication.mount()); err != nil {
		log.Printf("server failed to start %s", err)
		os.Exit(1)
	}
}
