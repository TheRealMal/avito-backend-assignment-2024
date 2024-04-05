package api

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
func (s Service) HandleUserBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
		return
	case len(feature_id) == 0:
		return
	}

	params := parseUserBannerParams(tag_id, feature_id, use_last_revision)
	if params == nil {
		return
	}
	w.Write([]byte("OK"))
}

func parseUserBannerParams(params ...string) *UserBannerParams {
	result := &UserBannerParams{}
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
