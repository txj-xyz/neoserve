package routes

// func handleUpload(w http.ResponseWriter, r *http.Request) {

// 	// Check for POST method
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Authenticate request
// 	authHeader := r.Header.Get("Authorization")
// 	if !strings.HasPrefix(authHeader, "Bearer ") || strings.TrimPrefix(authHeader, "Bearer ") != cfg.Auth.Token {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Parse the uploaded file
// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, "Failed to read file from form", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Generate a random file name
// 	randomName := config.GenerateRandomName() + regexp.MustCompile(`\.[a-zA-Z0-9]+$`).FindString(handler.Filename)

// 	logger.NewLogger(logger.DEBUG).Debug("Upload: [IP=%s] [name=%s] [random=%s] [uri=%s]", r.RemoteAddr, handler.Filename, randomName, r.RequestURI)

// 	// Save the file to the upload directory
// 	filePath := filepath.Join(cfg.Paths.Uploads, randomName)
// 	destFile, err := os.Create(filePath)
// 	if err != nil {
// 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// 		return
// 	}
// 	defer destFile.Close()

// 	if _, err := io.Copy(destFile, file); err != nil {
// 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "File uploaded successfully: %s\n", randomName)
// }