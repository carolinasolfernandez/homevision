package house

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetHouses(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a mock JSON response
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"houses":[{"id":1,"address":"123 Main St","homeowner":"John Doe","price":100000,"photoURL":"http://example.com/photo.jpg"}]}`)
	}))
	defer server.Close()

	// Initialize houseService
	hs := houseService{
		httpClient: server.Client(),
		url:        server.URL,
	}

	houseCh := make(chan []House)
	errorCh := make(chan error)

	// Call GetHouses
	go hs.GetHouses(1, 10, houseCh, errorCh)

	// Receive the result from the channel
	houses := <-houseCh

	// Check if the houses were received
	if len(houses) != 1 {
		t.Errorf("expected 1 house, got %d", len(houses))
	}
}

func TestSavePhotos(t *testing.T) {
	// Create a temporary directory for photos
	tmpDir, err := ioutil.TempDir("", "test_photos")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a mock photo content
		w.Header().Set("Content-Type", "image/jpeg")
		fmt.Fprint(w, "test photo content")
	}))
	defer server.Close()

	// Initialize houseService
	hs := houseService{
		httpClient: server.Client(),
		photosDir:  tmpDir,
	}

	// Create a channel for houseCh, errorCh, and doneCh
	houseCh := make(chan []House)
	errorCh := make(chan error)
	doneCh := make(chan struct{})

	// Send a test house to the house channel
	go func() {
		houseCh <- []House{{
			Id:        1,
			Address:   "123 Main St",
			Homeowner: "John Doe",
			Price:     100000,
			PhotoURL:  server.URL,
		}}
		close(houseCh)
	}()

	// Call SavePhotos
	go hs.SavePhotos(houseCh, doneCh, errorCh)

	// Wait for SavePhotos to finish
	<-doneCh

	// Get the host and port from the listener address
	addr := server.Listener.Addr().String()
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		t.Fatalf("unexpected listener address format: %s", addr)
	}
	port := parts[1]

	// TODO the extension should only be an image extension
	// Check if the photo file was saved
	filePath := tmpDir + "/1-123 Main St." + ".1:" + port
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("expected photo file to exist, got error:", err)
	}
}

func TestGetPage(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a mock JSON response
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"houses":[{"id":1,"address":"123 Main St","homeowner":"John Doe","price":100000,"photoURL":"http://example.com/photo.jpg"}]}`)
	}))
	defer server.Close()

	// Initialize houseService
	hs := houseService{
		httpClient: server.Client(),
		url:        server.URL,
	}

	housesCh := make(chan []House, 1)

	// Call getPage
	err := hs.getPage(1, 10, housesCh)

	// Check for errors
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Receive the result from the channel
	houses := <-housesCh

	// Check if the houses were received
	if len(houses) != 1 {
		t.Errorf("expected 1 house, got %d", len(houses))
	}
}

func TestSavePhoto(t *testing.T) {
	// Create a temporary directory for photos
	tmpDir, err := ioutil.TempDir("", "test_photos")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a mock photo content
		w.Header().Set("Content-Type", "image/jpeg")
		fmt.Fprint(w, "test photo content")
	}))
	defer server.Close()

	// Initialize houseService
	hs := houseService{
		httpClient: server.Client(),
		photosDir:  tmpDir,
	}

	// Call savePhoto
	err = hs.savePhoto("1-123 Main St.jpg", server.URL)

	// Check for errors
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the photo file was saved
	filePath := tmpDir + "/1-123 Main St.jpg"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("expected photo file to exist, got error:", err)
	}
}
