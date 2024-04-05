package server

import (
	handlers "avito-backend/pkg/api"
	"net/http"
)

func InitServerMux() *http.ServeMux {
	service := handlers.Service{}

	mux := http.NewServeMux()

	mux.Handle("/banner", handlers.AdminMiddleware(http.HandlerFunc(service.HandleBanner)))
	mux.Handle("/banner/", handlers.AdminMiddleware(http.HandlerFunc(service.HandleBannerID)))
	mux.Handle("/user_banner", handlers.AdminMiddleware(http.HandlerFunc(service.HandleUserBanner)))

	return mux
}
