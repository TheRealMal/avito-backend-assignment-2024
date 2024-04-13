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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const gracefulShutdownTime = 2 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to start logger: %s", err.Error())
	}
	defer logger.Sync() //nolint:errcheck

	// Load environment variables
	err = godotenv.Load(".env")
	if err != nil {
		logger.Fatal("failed to load .env file",
			zap.Error(err),
		)
	}
	dbHost, dbName, dbUser, dbPass := os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")

	// Connect to db
	logger.Info("connecting to db")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPass, dbHost, dbName)
	dbPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		logger.Fatal("failed to connect to db",
			zap.Error(err),
		)
	}
	defer dbPool.Close()
	database := db.InitDB(ctx, dbPool)

	// Init server mux
	s := handlers.NewServiceHandler(database, logger)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/banner", s.HandleBanner)
	adminMux.HandleFunc("/banner/", s.HandleBannerID)
	adminHandler := s.AdminMiddleware(adminMux)

	userMux := http.NewServeMux()
	userMux.HandleFunc("/user_banner", s.HandleUserBanner)
	userHandler := s.UserMiddleware(userMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/banner", adminHandler)
	siteMux.Handle("/user_banner", userHandler)

	httpHandler := s.LogMiddleware(siteMux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}

	// Start server
	go func() {
		logger.Info("starting server")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server error",
				zap.Error(err),
			)
		}
		logger.Info("http server stopped")
	}()

	// Wait for kill signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info("got termination signal")

	// Stop invalidator, graceful shutdown server
	// db connection is closed via defer
	logger.Info("gracefully shutting down")
	cancel()
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), gracefulShutdownTime)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("http server shutdown error",
			zap.Error(err),
		)
	}
	logger.Info("graceful shutdown complete")
}
