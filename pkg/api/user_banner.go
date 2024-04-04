package handlers

import "net/http"

// Get user banner request handler
func HandleUserBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	tag_id := r.URL.Query().Get("tag_id")
	feature_id := r.URL.Query().Get("feature_id")
	use_last_revision := r.URL.Query().Get("use_last_revision")

	if len(tag_id) == 0 || len(feature_id) == 0 {
		// Return error
	}

	// Interact with database

	w.Write([]byte("OK"))
}
