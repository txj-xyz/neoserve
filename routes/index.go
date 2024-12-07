package routes

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/txj-xyz/neoserve/file-server/internal/config"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.GetConfig()
	if err != nil {
		var r = fmt.Sprintf("failed to get any loaded config %v", err)
		http.Error(w, r, http.StatusInternalServerError)
	}

	http.ServeFile(w, r, filepath.Join(cfg.Paths.PublicPath(), "index.html"))
}
