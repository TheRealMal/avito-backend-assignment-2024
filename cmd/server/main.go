package main

import (
	"log"
	"net/http"

	api "avito-backend/pkg/api"
)

func main() {
	// Create generated server.
	srv, err := api.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
