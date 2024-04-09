package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

type BannerPatchBody struct {
	TagIDs    *[]int           `json:"tag_ids"`
	FeatureID *int             `json:"feature_id"`
	Content   *json.RawMessage `json:"content"`
	IsActive  *bool            `json:"is_active"`
}

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
	id, err := getBannerID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var patchParams BannerPatchBody
	err = decoder.Decode(&patchParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := s.db.UpdateBanner(id, patchParams.TagIDs, patchParams.FeatureID, patchParams.Content, patchParams.IsActive)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s ServiceHandler) HandleBannerIDDelete(w http.ResponseWriter, r *http.Request) {
	id, err := getBannerID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := s.db.DeleteBanner(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getBannerID(path string) (int, error) {
	str_id := strings.TrimPrefix(path, "/banner/")
	return strconv.Atoi(str_id)
}
