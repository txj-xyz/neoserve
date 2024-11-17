package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/txj-xyz/neoserve/file-server/internal/config"
	"github.com/txj-xyz/neoserve/file-server/internal/logger"
)

var (
	cfg *config.Config
	log *logger.Logger
)

func main() {
	// Load configuration
	var err error
	cfg, err = config.LoadConfig("./config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log = logger.NewLogger(cfg.Logging.Level)

	// Ensure directories exist
	os.MkdirAll(cfg.Paths.Uploads, os.ModePerm)

	http.HandleFunc("/", serveIndex)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.Paths.Uploads))))
	http.HandleFunc("/upload", handleUpload)

	fmt.Printf("Start making web requests with the following auth key: [%s]\n", cfg.Auth.Token)
	fmt.Printf("Server running on http://localhost:%s\n", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

// Server main root page
func serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Logf("[GET] [%s] [%s]", r.Host, r.RequestURI)
	http.ServeFile(w, r, filepath.Join(cfg.Paths.Public, "index.html"))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// Check for POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate request
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") || strings.TrimPrefix(authHeader, "Bearer ") != cfg.Auth.Token {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Debug("Found file: '%s' to upload.", handler.Filename)

	// Generate a random file name
	randomName := config.GenerateRandomName()

	// Save the file to the upload directory
	filePath := filepath.Join(cfg.Paths.Uploads, randomName)
	destFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", randomName)
}
