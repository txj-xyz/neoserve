package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/txj-xyz/neoserve/file-server/internal/config"
	"github.com/txj-xyz/neoserve/file-server/internal/logger"
	"github.com/txj-xyz/neoserve/file-server/routes"
)

var (
	cfg *config.Config
	log *logger.Logger
)

// 	http.HandleFunc("/", serveIndex)
// 	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.Paths.Uploads))))
// 	http.HandleFunc("/upload", handleUpload)

func main() {
	// Load config
	tmpCfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(1)
	}
	cfg = tmpCfg

	// Create a logger
	log = logger.NewLogger(cfg.Logging.Level)

	_, err = os.ReadDir(cfg.Paths.Uploads)
	if err != nil {
		log.Error("Error, cannot find uploads directory, creating it now")
		os.MkdirAll(cfg.Paths.Uploads, os.ModePerm)
	}

	// Create a Router for everything
	r := chi.NewRouter()

	// Setup middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Server timeout
	r.Use(middleware.Timeout(60 * time.Second))

	// ------------------------ Routes ------------------------
	// GET Routes
	r.Get("/", serveIndex)
	routes.FileListing(r, "/v1/files", http.Dir(cfg.Paths.UploadsPath())) // View the uploads

	// POST Routes
	r.Post("/v1/upload", routes.UploadFile)
	r.Post("/api/webhooks/{channelID}/{webhookToken}", routes.DiscordWebhookPassthrough)

	// Start Server

	if cfg.Server.DevMode {
		fmt.Printf("Neoserve DEV listening on %s\n", "http://localhost:8081")
		err = http.ListenAndServe(":8081", r)
		if err != nil {
			fmt.Printf("Error starting dev neoserver: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Neoserve listening on %s\n", cfg.Server.GenerateURL())
		err = http.ListenAndServe(":"+cfg.Server.Port, r)
		if err != nil {
			fmt.Printf("Error starting neoserver: %s\n", err)
			os.Exit(1)
		}
	}
	
}

// Server main root page
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(cfg.Paths.PublicPath(), "index.html"))
}
