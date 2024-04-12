package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRetryableTransport_RoundTrip(t *testing.T) {
	// Create a test server
	server := createTestServer()
	defer server.Close()

	// Create a new client with retry set to 3
	client := NewClient(3)

	// Create a request
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestExponentialBackoff(t *testing.T) {
	expectedBackoffs := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
	}

	for i, expected := range expectedBackoffs {
		backoff := backoff(i)
		if backoff != expected {
			t.Errorf("expected backoff %v, got %v", expected, backoff)
		}
	}
}

func TestRetryableStatusCodes(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{http.StatusInternalServerError, true},
		{http.StatusBadGateway, true},
		{http.StatusServiceUnavailable, true},
		{http.StatusGatewayTimeout, true},
		{http.StatusOK, false},
		{http.StatusNotFound, false},
	}

	for _, tc := range testCases {
		resp := &http.Response{StatusCode: tc.statusCode}

		retry := shouldRetry(nil, resp)
		if retry != tc.expected {
			t.Errorf("expected retry=%t for status code %d, got retry=%t", tc.expected, tc.statusCode, retry)
		}
	}
}

func createTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response
		w.WriteHeader(http.StatusOK)
	}))
}
