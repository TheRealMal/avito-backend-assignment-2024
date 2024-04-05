package server

import (
	"avito-backend/pkg/db"
	"avito-backend/pkg/handlers"
	"log"
	"net/http"
)

func InitServerMux(databaseURL string) *http.ServeMux {
	database, err := db.InitDB(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %s\n", err.Error())
	}

	s := handlers.NewServiceHandler(database)
	mux := http.NewServeMux()

	mux.Handle("/banner", handlers.AdminMiddleware(http.HandlerFunc(s.HandleBanner)))
	mux.Handle("/banner/", handlers.AdminMiddleware(http.HandlerFunc(s.HandleBannerID)))
	mux.Handle("/user_banner", handlers.UserMiddleware(http.HandlerFunc(s.HandleUserBanner)))

	return mux
}
