package api

import (
	"net/http"
	"strconv"
)

func (s Service) HandleBanner(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.HandleBannerGet(w, r)
	case http.MethodPost:
		s.HandleBannerPost(w, r)
	default:
		// return 500
	}
}

type BannerGetParams struct {
	TagID     int
	FeatureID int
	Limit     int
	Offset    int
}

func (s Service) HandleBannerGet(w http.ResponseWriter, r *http.Request) {
	tag_id := r.URL.Query().Get("tag_id")
	feature_id := r.URL.Query().Get("feature_id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Parse params to struct
	params := parseBannerGetParams(tag_id, feature_id, limit, offset)
	if params == nil {
		return
	}

	w.Write([]byte("OK"))
}

func parseBannerGetParams(params ...string) *BannerGetParams {
	result := &BannerGetParams{}
	var err error
	for idx, param := range params {
		switch idx {
		case 0:
			result.TagID, err = strconv.Atoi(param)
			if err != nil {
				return nil
			}
		case 1:
			result.FeatureID, err = strconv.Atoi(param)
			if err != nil {
				return nil
			}
		case 2:
			result.Limit, err = strconv.Atoi(param)
			if err != nil {
				return nil
			}
		case 3:
			result.Offset, err = strconv.Atoi(param)
			if err != nil {
				return nil
			}
		}
	}
	return result
}

func (s Service) HandleBannerPost(w http.ResponseWriter, r *http.Request) {
	// JSON Request Body
}
