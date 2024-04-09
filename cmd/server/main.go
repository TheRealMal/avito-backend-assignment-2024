package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avito-backend/internal/db"
	"avito-backend/internal/handlers"

	"github.com/joho/godotenv"
)

const gracefulShutdownTime = 2 * time.Second

func main() {
	// Load environment variables
	_ = godotenv.Load(".env")
	dbHost, dbName, dbUser, dbPass := os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")

	// Connect to db
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPass, dbHost, dbName)
	database, err := db.InitDB(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %s\n", err.Error())
	}
	defer database.Close()

	// Start db cache invalidator
	ctx, cancel := context.WithCancel(context.Background())
	go database.StartInvalidator(ctx)

	// Init server mux
	s := handlers.NewServiceHandler(database)
	mux := http.NewServeMux()
	mux.Handle("/banner", handlers.AdminMiddleware(http.HandlerFunc(s.HandleBanner)))
	mux.Handle("/banner/", handlers.AdminMiddleware(http.HandlerFunc(s.HandleBannerID)))
	mux.Handle("/user_banner", handlers.UserMiddleware(http.HandlerFunc(s.HandleUserBanner)))
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
		log.Println("http server stopped")
	}()

	// Wait for kill signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Stop invalidator, graceful shutdown server
	// db connection is closed via defer
	cancel()
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), gracefulShutdownTime)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("http shutdown error: %v", err)
	}
	log.Println("graceful shutdown complete.")
}
