package routes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/txj-xyz/neoserve/file-server/internal/config"
	"github.com/txj-xyz/neoserve/file-server/internal/webhook"
)

// Middleware for IP whitelist
func AdminIPWhitelist(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := config.GetConfig()
		ipStr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ipStr = r.RemoteAddr
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			http.Error(w, "Forbidden (Invalid IP)", http.StatusForbidden)
			return
		}

		allowed := false
		for _, allowedIP := range cfg.Admin.IPWhitelist {
			if strings.Contains(allowedIP, "/") {
				_, ipnet, err := net.ParseCIDR(allowedIP)
				if err == nil && ipnet.Contains(ip) {
					allowed = true
					break
				}
			} else {
				allowedIPParsed := net.ParseIP(allowedIP)
				if allowedIPParsed != nil && allowedIPParsed.Equal(ip) {
					allowed = true
					break
				}
			}
		}
		if !allowed {
			http.Error(w, "Forbidden (IP)", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var (
	adminSessions   = make(map[string]struct{})
	adminSessionsMu sync.Mutex
)

// Session check (secure random token)
func AdminSessionRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || !isValidAdminSession(cookie.Value) {
			accept := r.Header.Get("Accept")
			xrw := r.Header.Get("X-Requested-With")
			if strings.Contains(accept, "application/json") || xrw == "XMLHttpRequest" || r.URL.Path == "/admin/session-check" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"status":"ERROR","msg":"Unauthorized"}`))
			} else {
				http.Redirect(w, r, "/admin/login", http.StatusFound)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isValidAdminSession(token string) bool {
	// Log session validation for dev_mode

	adminSessionsMu.Lock()
	defer adminSessionsMu.Unlock()
	_, ok := adminSessions[token]
	return ok
}

// Serve admin panel HTML
func ServeAdminPanel(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/admin.html")
}

// Admin login handler
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.GetConfig()
	type creds struct{ Username, Password string }
	var c creds
	json.NewDecoder(r.Body).Decode(&c)
	if c.Username == cfg.Admin.Username && c.Password == cfg.Admin.Password {
		token := generateSessionToken()
		adminSessionsMu.Lock()
		adminSessions[token] = struct{}{}
		adminSessionsMu.Unlock()
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_session",
			Value:    token,
			Path:     "/",
			Domain:   cfg.Server.Domain,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   86400, // 24 hours
		})
		w.Write([]byte(`{"status":"OK"}`))
	} else {
		w.Write([]byte(`{"status":"ERROR","msg":"Invalid credentials"}`))
	}
}

func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Admin logout handler
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin_session")
	if err == nil {
		adminSessionsMu.Lock()
		delete(adminSessions, cookie.Value)
		adminSessionsMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: "admin_session", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	w.Write([]byte(`{"status":"OK"}`))
}

// Handle admin file upload
func AdminUploadFile(w http.ResponseWriter, r *http.Request) {

	cfg, err := config.GetConfig()
	if err != nil {
		http.Error(w, "Internal Error Occurred", http.StatusInternalServerError)
		return
	}

	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, cfg.Server.MaxUpload*1024*1024)
	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"ERROR","msg":"Invalid upload"}`))
		return
	}
	defer uploadedFile.Close()

	// Generate random file name
	ext := ""
	if dot := strings.LastIndex(handler.Filename, "."); dot != -1 {
		ext = handler.Filename[dot:]
	}
	rnd := config.GenerateRandomName() + ext

	savePath := filepath.Join(cfg.Paths.Uploads, rnd)
	savedFile, err := os.Create(savePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","msg":"Unable to save file"}`))
		return
	}
	defer savedFile.Close()

	if _, err := io.Copy(savedFile, uploadedFile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","msg":"Unable to save file"}`))
		return
	}

	// Discord webhook log if enabled
	link := cfg.Server.GenerateURL() + "/v1/files/" + rnd
	if cfg.Logging.Discord.Enabled {
		imported_webhook := false
		webhook.SendWebhookLog(link, rnd)
		imported_webhook = true
		_ = imported_webhook
	}

	w.Write([]byte(`{"status":"OK","url":"` + link + `"}`))
}

type FileListResponse struct {
	Files      []string `json:"files"`
	TotalCount int      `json:"totalCount"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
	TotalPages int      `json:"totalPages"`
}

// List uploaded files (JSON) with pagination support
func AdminListFiles(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.GetConfig()

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 50

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 200 {
			pageSize = ps
		}
	}

	// Collect all files
	var allFiles []string
	filepath.WalkDir(cfg.Paths.Uploads, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			allFiles = append(allFiles, d.Name())
		}
		return nil
	})

	sort.Strings(allFiles)
	totalCount := len(allFiles)
	totalPages := (totalCount + pageSize - 1) / pageSize

	// Adjust page if it's out of bounds
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// Calculate slice bounds
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalCount {
		end = totalCount
	}

	// Get the slice of files for the current page
	var files []string
	if start < totalCount {
		files = allFiles[start:end]
	} else {
		files = []string{}
	}

	// Create response
	response := FileListResponse{
		Files:      files,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Admin session check endpoint
func AdminSessionCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"OK"}`))
}

// Bulk delete files via admin session (no Bearer token required)
func AdminBulkDeleteFiles(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.GetConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","msg":"Internal Error Occurred"}`))
		return
	}
	var files []string
	if err := json.NewDecoder(r.Body).Decode(&files); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"ERROR","msg":"Invalid request body"}`))
		return
	}

	type FileDeleteResult struct {
		File   string `json:"file"`
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	results := make([]FileDeleteResult, 0, len(files))
	for _, filename := range files {
		if filename == "" {
			results = append(results, FileDeleteResult{File: filename, Status: "ERROR", Msg: "Invalid filename"})
			continue
		}
		filePath := filepath.Join(cfg.Paths.Uploads, filename)
		err := os.Remove(filePath)
		if err != nil {
			results = append(results, FileDeleteResult{File: filename, Status: "ERROR", Msg: "File not found or could not be deleted"})
		} else {
			results = append(results, FileDeleteResult{File: filename, Status: "OK", Msg: "File deleted successfully"})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Delete file via admin session (no Bearer token required)
func AdminDeleteFile(w http.ResponseWriter, r *http.Request) {

	cfg, err := config.GetConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","msg":"Internal Error Occurred"}`))
		return
	}

	filename := chi.URLParam(r, "filename")
	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"ERROR","msg":"Invalid filename"}`))
		return
	}

	filePath := filepath.Join(cfg.Paths.Uploads, filename)
	err = os.Remove(filePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"ERROR","msg":"File not found or could not be deleted"}`))
		return
	}

	w.Write([]byte(`{"status":"OK","msg":"File deleted successfully"}`))
}
