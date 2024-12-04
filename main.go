package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/txj-xyz/neoserve/file-server/internal/config"
	"github.com/txj-xyz/neoserve/file-server/internal/logger"
	"github.com/txj-xyz/neoserve/file-server/routes"
)

var (
	log *logger.Logger
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(1)
	}

	// Create a logger
	log = logger.New()

	// If the uploads directory doesnt exist then create it
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
	r.Use(middleware.Timeout(60 * time.Second)) // Set server max timeouts

	// GET Routes
	r.Get("/", routes.ServeIndex)
	routes.FileListing(r, "/v1/files", http.Dir(cfg.Paths.UploadsPath())) // View the uploads

	// POST Routes
	r.Post("/v1/upload", routes.UploadFile)
	r.Post("/api/webhooks/{channelID}/{webhookToken}", routes.DiscordWebhookPassthrough)

	// Start up neoserve with the router we have
	Neoserve(r, cfg)
}

// Main starting logic to check for config flags etc
func Neoserve(r *chi.Mux, cfg *config.Config) {

	log.Debug("loading server components..")

	var listenAddr string = ":" + cfg.Server.Port

	if cfg.Server.DevMode {
		log.Info(
			"neoserve starting [DEV MODE]",
			"port", 8081,
			"url", "http://localhost:8081",
		)
		listenAddr = ":8081"
	}

	log.Info(
		"neoserve starting",
		"port", cfg.Server.Port,
		"url", cfg.Server.GenerateURL(),
	)
	// Start up the web server
	err := http.ListenAndServe(listenAddr, r)
	if err != nil {
		panic(fmt.Errorf("failed to bind port for neoserve: %v", err))
	}
}
