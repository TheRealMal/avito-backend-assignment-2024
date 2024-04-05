package main

import (
	"log"
	"net/http"
	"os"

	server "avito-backend/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed loading .env: %s\n", err.Error())
	}

	srv := server.InitServerMux(os.Getenv("DB_URL"))
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("http server error: %s\n", err)
	}
}
