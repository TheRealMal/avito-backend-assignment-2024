package main

import (
	"log"
	"net/http"

	server "avito-backend/pkg/server"
)

func main() {
	srv := server.InitServerMux()
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
