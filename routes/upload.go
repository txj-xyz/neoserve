package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/txj-xyz/neoserve/file-server/internal/config"
	"github.com/txj-xyz/neoserve/file-server/internal/webhook"
)


// Neoserve file uploader thingy
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
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

	// Parse the uploaded file
	r.Body = http.MaxBytesReader(w, r.Body, cfg.Server.MaxUpload * 1024 * 1024) // 45Mb limit
	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
			res := UploadResponse{
			Status: "ERROR",
			Code:   400,
			URL:    "null",
		}
		http.Error(w, res.ToJSON(), http.StatusBadRequest)
		return
	}
	defer uploadedFile.Close()

	// Generate a random file name
	rnd := config.GenerateRandomName() + regexp.MustCompile(`\.[a-zA-Z0-9]+$`).FindString(handler.Filename)

	// Save the file
	savePath := filepath.Join(cfg.Paths.Uploads, rnd)
	savedFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer savedFile.Close()

	if _, err := io.Copy(savedFile, uploadedFile); err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Send out a discord webhook log if its enabled
	link := fmt.Sprintf("%s/v1/files/%s", cfg.Server.GenerateURL(), rnd)
	webhook.SendWebhookLog(link, rnd)

	// Done send back a response
	res := UploadResponse{
		Status: "OK",
		Code:   200,
		URL:    fmt.Sprintf("%s/v1/files/%s", cfg.Server.GenerateURL(), rnd),
	}
	w.Write([]byte(res.ToJSON()))

}
