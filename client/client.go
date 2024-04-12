package client

import (
	"bytes"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

type retryableTransport struct {
	transport http.RoundTripper
	retries   int
}

func NewClient(retries int) *http.Client {
	transport := &retryableTransport{
		transport: &http.Transport{},
		retries:   retries,
	}

	return &http.Client{
		Transport: transport,
	}
}

func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Send the request
	resp, err := t.transport.RoundTrip(req)

	// Retry logic
	retry := 0

	for shouldRetry(err, resp) && retry < t.retries {
		// Wait for the specified backoff period
		time.Sleep(backoff(retry))
		// We're going to retry, consume any response to reuse the connection.
		drainBody(resp)
		// Clone the request body again
		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Retry the request
		resp, err = t.transport.RoundTrip(req)
		retry++
	}
	return resp, err
}

func drainBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}

	if resp.StatusCode == http.StatusInternalServerError ||
		resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout {
		log.Println("retrying request")
		return true
	}
	return false
}

func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}
