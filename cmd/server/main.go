package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	server "avito-backend/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	dbHost, dbName, dbUser, dbPass := os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPass, dbHost, dbName)

	srv := server.InitServerMux(databaseURL)
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("http server error: %s\n", err)
	}
}
