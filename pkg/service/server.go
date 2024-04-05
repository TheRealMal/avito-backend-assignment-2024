package server

import (
	handlers "avito-backend/pkg/api"
	"net/http"
)

func InitServerMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/banner", handlers.AdminMiddleware(http.HandlerFunc(handlers.HandleBanner)))
	mux.Handle("/banner/", handlers.AdminMiddleware(http.HandlerFunc(handlers.HandleBannerID)))
	mux.Handle("/user_banner", handlers.AdminMiddleware(http.HandlerFunc(handlers.HandleUserBanner)))

	return mux
}
