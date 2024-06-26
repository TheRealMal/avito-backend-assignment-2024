package handlers

import (
	"avito-backend/internal/db"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
)

const (
	DefaultLimit  = 100
	DefaultOffset = 0
)

func (s ServiceHandler) HandleBanner(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.HandleBannerGet(w, r)
	case http.MethodPost:
		s.HandleBannerPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type BannerGetParams struct {
	TagID     int
	FeatureID int
	Limit     int
	Offset    int
}

func (s ServiceHandler) HandleBannerGet(w http.ResponseWriter, r *http.Request) {
	tagID := r.URL.Query().Get("tag_id")
	featureID := r.URL.Query().Get("feature_id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Parse params to struct
	params := parseBannerGetParams(tagID, featureID, limit, offset)
	if params == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	banners, err := s.db.GetBanners(params.TagID, params.FeatureID, params.Limit, params.Offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(banners)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parseBannerGetParams(params ...string) *BannerGetParams {
	result := &BannerGetParams{}
	var err error
	for idx, param := range params {
		switch idx {
		case 0:
			result.TagID, err = strconv.Atoi(param)
			if err != nil {
				result.TagID = -1
			}
		case 1:
			result.FeatureID, err = strconv.Atoi(param)
			if err != nil {
				result.FeatureID = -1
			}
		case 2:
			result.Limit, err = strconv.Atoi(param)
			if err != nil {
				result.Limit = DefaultLimit
			}
		case 3:
			result.Offset, err = strconv.Atoi(param)
			if err != nil {
				result.Offset = DefaultOffset
			}
		}
	}
	return result
}

func (s ServiceHandler) HandleBannerPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var banner db.Banner
	err := decoder.Decode(&banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.db.CreateBanner(banner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
