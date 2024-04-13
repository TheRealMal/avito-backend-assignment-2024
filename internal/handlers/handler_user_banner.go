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
	queryValues := r.URL.Query()
	tagID := queryValues.Get("tag_id")
	featureID := queryValues.Get("feature_id")
	useLastRevision := queryValues.Get("use_last_revision")

	// Check required params and some formats
	switch {
	case len(tagID) == 0:
		w.WriteHeader(http.StatusBadRequest)
		return
	case len(featureID) == 0:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := parseUserBannerParams(tagID, featureID, useLastRevision)
	if params == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, isActive, err := s.db.GetBannerContent(params.TagID, params.FeatureID, params.UseLastRevision)
	switch {
	case err == nil && response != nil:
		// Get role from ctx
		untypedRole := r.Context().Value(ContextRoleKey)
		if untypedRole == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Cast to string
		role, ok := untypedRole.(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// If not active and not admin
		if !isActive && role != AdminRole {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(*response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
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
