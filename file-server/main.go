package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/txj-xyz/neoserve/file-server/internal/config"
	"github.com/txj-xyz/neoserve/file-server/internal/logger"
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
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
	}
	
	// Create a logger
	log = logger.NewLogger(cfg.Logging.Level)

	// Show the dir layout in dev mode only
	if cfg.Server.DevMode {
		log.Debug("Printing out the dir tree for stack testing")
		config.PrintDirTree("./", "")
	}

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
	r.Use(middleware.URLFormat)

	// Server timeout
	r.Use(middleware.Timeout(60 * time.Second))

	// ------------------------ Routes ------------------------
	// GET Routes
	r.Get("/", serveIndex)
	// r.Get("/uploads")

	// POST Routes
	r.Post("/upload", handleUpload)


	// Start Server
	log.Logf("Serving on http://localhost:%s\n", cfg.Server.Port)
	http.ListenAndServe(":"+cfg.Server.Port, r)
}

// Server main root page
func serveIndex(w http.ResponseWriter, r *http.Request) {

	// log.Logf("Index file path: %s", filepath.Join(cfg.Paths.Public, "index.html"))
	// http.ServeFile(w, r, filepath.Join(cfg.Paths.Public, "index.html"))
	w.Write([]byte("root."))
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


	// Generate a random file name
	randomName := config.GenerateRandomName() + regexp.MustCompile(`\.[a-zA-Z0-9]+$`).FindString(handler.Filename)

	log.Debug("Upload: [IP=%s] [name=%s] [random=%s] [uri=%s]", r.RemoteAddr, handler.Filename, randomName, r.RequestURI)


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
