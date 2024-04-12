package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunctionality(t *testing.T) {
	// Create a temporary directory for photos
	tmpDir, err := ioutil.TempDir("", "test_photos")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test HTTP server for image data
	serverImage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a mock photo content
		w.Header().Set("Content-Type", "image/jpeg")
		fmt.Fprint(w, "test photo content")
	}))
	defer serverImage.Close()

	// Create a mock HTTP server for house data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with mock JSON data
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"houses":[{"id":1,"address":"123 Main St","homeowner":"John Doe","price":100000,"photoURL":"` + serverImage.URL + `"}]}`))
	}))

	// Set environment variables for test configuration
	os.Setenv("HOUSES_URL", server.URL)
	os.Setenv("NUM_PAGES", "1")
	os.Setenv("NUM_PER_PAGE", "1")
	os.Setenv("PHOTOS_DIR", tmpDir)
	os.Setenv("CLIENT_RETRIES", "3")

	defer server.Close()

	// Call the main function
	main()

	// Get the host and port from the listener address
	addr := serverImage.Listener.Addr().String()
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		t.Fatalf("unexpected listener address format: %s", addr)
	}
	port := parts[1]

	// TODO the extension should only be an image extension
	// Assert that the photo was saved
	filePath := os.Getenv("PHOTOS_DIR") + "/1-123 Main St." + ".1:" + port
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "photo file should be saved without error")
}
