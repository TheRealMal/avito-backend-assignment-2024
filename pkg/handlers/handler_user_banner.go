package handlers

import (
	"net/http"
	"strconv"
)

type UserBannerParams struct {
	TagID           int
	FeatureID       int
	UseLastRevision bool
}

// Get user banner request handler
func (s ServiceHandler) HandleUserBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Take params from URL
	query_values := r.URL.Query()
	tag_id := query_values.Get("tag_id")
	feature_id := query_values.Get("feature_id")
	use_last_revision := query_values.Get("use_last_revision")

	// Check required params and some formats
	switch {
	case len(tag_id) == 0:
		w.WriteHeader(http.StatusBadRequest)
		return
	case len(feature_id) == 0:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := parseUserBannerParams(tag_id, feature_id, use_last_revision)
	if params == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := s.db.GetBannerContent(params.TagID, params.FeatureID, params.UseLastRevision)
	switch {
	case err == nil && response != nil:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(*response)
	case err == nil && response == nil:
		w.WriteHeader(http.StatusNotFound)
	case err != nil && response == nil:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parseUserBannerParams(params ...string) *UserBannerParams {
	result := &UserBannerParams{
		UseLastRevision: false,
	}
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
			if len(param) != 0 && param != "false" && param != "true" {
				return nil
			}
			result.UseLastRevision = param == "true"
		}
	}
	return result
}
