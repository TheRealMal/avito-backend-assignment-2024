package api

import (
	"net/http"
	"strconv"
	"strings"
)

func (s Service) HandleBannerID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		s.HandleBannerIDPatch(w, r)
	case http.MethodDelete:
		s.HandleBannerIDDelete(w, r)
	default:
		// return 500
	}
}

func (s Service) HandleBannerIDPatch(w http.ResponseWriter, r *http.Request) {
	id, err := strings.TrimPrefix(r.URL.Path, "/banner/")
	if err != nil {
		return
	}
}

func (s Service) HandleBannerIDDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strings.TrimPrefix(r.URL.Path, "/banner/")
	if err != nil {
		return
	}
}

func getBannerID(path string) (int, error) {
	str_id := strings.TrimPrefix(path, "/banner/")
	return strconv.Atoi(str_id)
}
