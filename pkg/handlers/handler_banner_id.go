package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func (s ServiceHandler) HandleBannerID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		s.HandleBannerIDPatch(w, r)
	case http.MethodDelete:
		s.HandleBannerIDDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (s ServiceHandler) HandleBannerIDPatch(w http.ResponseWriter, r *http.Request) {
	_, err := getBannerID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s ServiceHandler) HandleBannerIDDelete(w http.ResponseWriter, r *http.Request) {
	_, err := getBannerID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getBannerID(path string) (int, error) {
	str_id := strings.TrimPrefix(path, "/banner/")
	return strconv.Atoi(str_id)
}
