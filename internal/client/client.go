package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/I-invincib1e/http-cli/internal/config"
)

// Response holds the HTTP response data
type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
}

// ExecuteRequest executes an HTTP request based on the provided config
func ExecuteRequest(cfg *config.Config) (*Response, error) {
	// Add bearer token to headers if provided
	if cfg.BearerToken != "" {
		if cfg.Headers == nil {
			cfg.Headers = make(map[string]string)
		}
		cfg.Headers["Authorization"] = "Bearer " + cfg.BearerToken
	}

	// Add basic auth if provided
	if cfg.BasicAuth != "" {
		parts := strings.SplitN(cfg.BasicAuth, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("basic auth format should be 'user:pass'")
		}
		if cfg.Headers == nil {
			cfg.Headers = make(map[string]string)
		}
		auth := base64.StdEncoding.EncodeToString([]byte(cfg.BasicAuth))
		cfg.Headers["Authorization"] = "Basic " + auth
	}

	// Create HTTP client
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	// Configure redirect following
	if !cfg.FollowRedirects {
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Create request
	var req *http.Request
	var err error

	if cfg.Body != "" {
		req, err = http.NewRequest(cfg.Method, cfg.URL, bytes.NewBufferString(cfg.Body))
	} else {
		req, err = http.NewRequest(cfg.Method, cfg.URL, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	// If body exists and no Content-Type header, default to application/json
	if cfg.Body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make request
	start := time.Now()
	resp, err := httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       bodyBytes,
		Duration:   duration,
	}, nil
}

