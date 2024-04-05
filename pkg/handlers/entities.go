package handlers

import (
	"avito-backend/pkg/db"
	"net/http"
)

type ServiceHandler struct {
	db db.Database
}

func NewServiceHandler(database db.Database) ServiceHandler {
	return ServiceHandler{
		db: database,
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
