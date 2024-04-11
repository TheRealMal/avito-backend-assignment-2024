package handlers

import (
	"avito-backend/internal/db"
	"net/http"

	"go.uber.org/zap"
)

type ServiceHandler struct {
	db  db.Database
	log *zap.Logger
}

func NewServiceHandler(database db.Database, logger *zap.Logger) ServiceHandler {
	return ServiceHandler{
		db:  database,
		log: logger,
	}
}

type Service interface {
	HandleBannerID(http.ResponseWriter, *http.Request)
	HandleBannerIDPatch(http.ResponseWriter, *http.Request)
	HandleBannerIDDelete(http.ResponseWriter, *http.Request)
	HandleBanner(http.ResponseWriter, *http.Request)
	HandleBannerGet(http.ResponseWriter, *http.Request)
	HandleBannerPost(http.ResponseWriter, *http.Request)
	HandleUserBanner(http.ResponseWriter, *http.Request)
}
