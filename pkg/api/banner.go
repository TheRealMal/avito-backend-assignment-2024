package handlers

import "net/http"

func HandleBanner(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandleBannerGet(w, r)
	case http.MethodPost:
		HandleBannerPost(w, r)
	default:
		// Return Error
	}
}

func HandleBannerGet(w http.ResponseWriter, r *http.Request) {
	tag_id := r.URL.Query().Get("tag_id")
	feature_id := r.URL.Query().Get("feature_id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

}

func HandleBannerPost(w http.ResponseWriter, r *http.Request) {
	// JSON Request Body
}
