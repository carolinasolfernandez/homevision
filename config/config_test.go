package config

import (
	"os"
	"testing"
)

func TestLoadEnvVariables(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("HOUSES_URL", "http://example.com")
	os.Setenv("NUM_PAGES", "10")
	os.Setenv("NUM_PER_PAGE", "20")
	os.Setenv("PHOTOS_DIR", "/path/to/photos")
	os.Setenv("CLIENT_RETRIES", "3")

	defer func() {
		// Clean up environment variables after testing
		os.Clearenv()
	}()

	// Load configuration
	cfg := loadEnvVariables(".")

	// Check if configuration values are loaded correctly
	if cfg.HousesUrl != "http://example.com" {
		t.Errorf("expected HOUSES_URL=%s, got %s", "http://example.com", cfg.HousesUrl)
	}
	if cfg.NumPages != 10 {
		t.Errorf("expected NUM_PAGES=%d, got %d", 10, cfg.NumPages)
	}
	if cfg.NumPerPage != 20 {
		t.Errorf("expected NUM_PER_PAGE=%d, got %d", 20, cfg.NumPerPage)
	}
	if cfg.PhotosDir != "/path/to/photos" {
		t.Errorf("expected PHOTOS_DIR=%s, got %s", "/path/to/photos", cfg.PhotosDir)
	}
	if cfg.ClientRetries != 3 {
		t.Errorf("expected CLIENT_RETRIES=%d, got %d", 3, cfg.ClientRetries)
	}
}
