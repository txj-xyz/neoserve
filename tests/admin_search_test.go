package routes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/txj-xyz/neoserve/file-server/internal/config"
)

// TestAdminFileSearch tests that searching for 'file117' returns expected files
func TestAdminFileSearch(t *testing.T) {
	// Load config from config.example.yaml (since config.yaml was not found)
	cfg, err := config.LoadConfig(filepath.Join("..", "config.yaml"))
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	uploadsDir := cfg.Paths.Uploads
	if uploadsDir == "" {
		t.Fatal("Uploads directory not set in config")
	}

	var foundFiles []string
	err = filepath.Walk(uploadsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(info.Name(), "file117.txt") {
			foundFiles = append(foundFiles, info.Name())
		}
		return nil
	})
	if err != nil {
		t.Fatalf("error walking uploads directory: %v", err)
	}

	if len(foundFiles) == 0 {
		t.Errorf("No files found matching 'file117'")
	}

	// Optionally, print the found files for debugging
	for _, f := range foundFiles {
		t.Logf("Found: %s", f)
	}
}
