package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gophish/gophish/models"
)

type shortenRequest struct {
	URL        string `json:"url"`
	CampaignId int64  `json:"campaign_id"`
	RId        string `json:"rid"`
}

// ShortenURL creates a built-in short link and returns the short URL.
// POST /api/shorten  {"url": "...", "campaign_id": 1, "rid": "AbCdEfG"}
// Returns {"short_url": "http://host/r/AbCdEfG"}
func (as *Server) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	sl, err := models.CreateShortLink(req.URL, req.CampaignId, req.RId)
	if err != nil {
		http.Error(w, `{"error":"failed to create short link"}`, http.StatusInternalServerError)
		return
	}

	// Build the short URL using the same host as the request (admin server).
	// The phish server is on a different port; we extract the base from the
	// original URL instead so the redirect points to the right server.
	base := extractBase(req.URL)
	shortURL := base + "/r/" + sl.Code

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": shortURL,
		"code":      sl.Code,
	})
}

// ExtractBase returns scheme+host from a full URL string.
func ExtractBase(rawURL string) string {
	return extractBase(rawURL)
}

// extractBase returns scheme+host from a full URL string.
func extractBase(rawURL string) string {
	// Find "://"
	idx := strings.Index(rawURL, "://")
	if idx < 0 {
		return ""
	}
	rest := rawURL[idx+3:]
	// Find next slash
	slash := strings.Index(rest, "/")
	if slash < 0 {
		return rawURL
	}
	return rawURL[:idx+3+slash]
}
