package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/otiai10/gosseract/v2"
)

const (
	defaultPort     = "8080"
	defaultLanguage = "eng"
	defaultTimeout  = 30 // seconds
	defaultMaxSize  = 10 // MB
)

// Supported languages map
var supportedLanguages = map[string]string{
	"english": "eng",
	"german":  "deu",
	"french":  "fra",
	"italian": "ita",
	// Also support direct tesseract language codes
	"eng": "eng",
	"deu": "deu",
	"fra": "fra",
	"ita": "ita",
}

// Configuration holds server settings
type Configuration struct {
	Port        string
	Timeout     int // seconds
	MaxFileSize int // MB
}

func main() {
	// Load configuration from environment variables
	config := loadConfiguration()

	// Create a new server with timeout configuration
	server := &http.Server{
		Addr:         ":" + config.Port,
		ReadTimeout:  time.Duration(config.Timeout) * time.Second,
		WriteTimeout: time.Duration(config.Timeout) * time.Second,
		IdleTimeout:  time.Duration(config.Timeout) * time.Second,
	}

	// Register the OCR handler
	http.HandleFunc("POST /ocr", func(w http.ResponseWriter, r *http.Request) {
		ocrHandler(w, r, config.MaxFileSize)
	})

	// Register health endpoint
	http.HandleFunc("GET /health", healthHandler)

	// Start the server
	log.Printf("Starting OCR server on port %s", config.Port)
	log.Printf("Timeout: %d seconds", config.Timeout)
	log.Printf("Max file size: %d MB", config.MaxFileSize)
	log.Printf("Supported languages: English (eng), German (deu), French (fra), Italian (ita)")
	log.Fatal(server.ListenAndServe())
}

// loadConfiguration loads configuration from environment variables
func loadConfiguration() Configuration {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Get timeout from environment variable or use default
	timeout := defaultTimeout
	if timeoutStr := os.Getenv("OCR_TIMEOUT"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = t
		} else {
			log.Printf("Invalid timeout value: %s, using default: %d seconds", timeoutStr, defaultTimeout)
		}
	}

	// Get max file size from environment variable or use default
	maxSize := defaultMaxSize
	if maxSizeStr := os.Getenv("OCR_MAX_FILE_SIZE"); maxSizeStr != "" {
		if size, err := strconv.Atoi(maxSizeStr); err == nil {
			maxSize = size
		} else {
			log.Printf("Invalid max file size value: %s, using default: %d MB", maxSizeStr, defaultMaxSize)
		}
	}

	return Configuration{
		Port:        port,
		Timeout:     timeout,
		MaxFileSize: maxSize,
	}
}

func ocrHandler(w http.ResponseWriter, r *http.Request, maxSizeMB int) {
	// Parse the multipart form with configurable max size
	err := r.ParseMultipartForm(int64(maxSizeMB) << 20) // Convert MB to bytes
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Get language parameter (default to English if not provided)
	lang := strings.ToLower(r.FormValue("lang"))
	tesseractLang := defaultLanguage
	if lang != "" {
		if code, exists := supportedLanguages[lang]; exists {
			tesseractLang = code
		} else {
			http.Error(w, "Unsupported language. Supported languages: english, german, french, italian", http.StatusBadRequest)
			return
		}
	}

	// Initialize new gosseract client
	client := gosseract.NewClient()
	defer client.Close()

	// Set language
	err = client.SetLanguage(tesseractLang)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to set language: %v", err), http.StatusInternalServerError)
		return
	}

	// Set image from bytes
	err = client.SetImageFromBytes(fileBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load image: %v", err), http.StatusInternalServerError)
		return
	}

	// Extract text
	text, err := client.Text()
	if err != nil {
		http.Error(w, fmt.Sprintf("OCR processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the OCR text
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}

// healthHandler provides a simple health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
