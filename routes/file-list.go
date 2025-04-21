package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/txj-xyz/neoserve/file-server/internal/config"
)

type UploadResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	URL    string `json:"url"`
}

func (r *UploadResponse) ToJSON() string {
	prettyResponse, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return string("")
	}
	return string(prettyResponse)
}

// File listing page
func FileListing(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		panic("neoserver does not permit any URL params.")
	}

	// Handle the path and redirect if necessary
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusPermanentRedirect).ServeHTTP)
		path += "/"
	}

	// Restrict the /v1/files/ directory itself (authentication)
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		// Authenticate request for the directory itself
		authHeader := r.Header.Get("Authorization")
		cfg, err := config.GetConfig()
		if err != nil {
			http.Error(w, "Internal Error Occurred", http.StatusInternalServerError)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") || strings.TrimPrefix(authHeader, "Bearer ") != cfg.Auth.Token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If authenticated, redirect to the file listing (or serve file list page)
		http.Redirect(w, r, path+"/", http.StatusFound)
	})

	// Serve files under /v1/files/ path without authentication (files themselves are not restricted)
	r.Get(path+"*", func(w http.ResponseWriter, r *http.Request) {
		// Serve the file if no authentication is required for this file path
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}