package todo

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer() {
	fmt.Println("Starting server...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
