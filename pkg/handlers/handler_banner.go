package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	DEFAULT_LIMIT  = 100
	DEFAULT_OFFSET = 0
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
	tag_id := r.URL.Query().Get("tag_id")
	feature_id := r.URL.Query().Get("feature_id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Parse params to struct
	params := parseBannerGetParams(tag_id, feature_id, limit, offset)
	if params == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(params)
	banners, err := s.db.GetBanners(params.TagID, params.FeatureID, params.Limit, params.Offset)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(banners)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
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
				result.Limit = DEFAULT_LIMIT
			}
		case 3:
			result.Offset, err = strconv.Atoi(param)
			if err != nil {
				result.Offset = DEFAULT_OFFSET
			}
		}
	}
	return result
}

func (s ServiceHandler) HandleBannerPost(w http.ResponseWriter, r *http.Request) {
	// JSON Request Body
}
