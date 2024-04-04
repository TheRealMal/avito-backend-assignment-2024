package handlers

import (
	"net/http"
	"strings"
)

func HandleBannerID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		HandleBannerIDPatch(w, r)
	case http.MethodDelete:
		HandleBannerIDDelete(w, r)
	default:
		// return 500
	}
}

func HandleBannerIDPatch(w http.ResponseWriter, r *http.Request) {
	str_id := strings.TrimPrefix(r.URL.Path, "/banner/")
}

func HandleBannerIDDelete(w http.ResponseWriter, r *http.Request) {
	str_id := strings.TrimPrefix(r.URL.Path, "/banner/")
}
