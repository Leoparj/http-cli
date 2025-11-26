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

// addAuth adds authentication headers to the config
func addAuth(cfg *config.Config) error {
    if cfg.Headers == nil {
        cfg.Headers = make(map[string]string)
    }

    if cfg.BearerToken != "" {
        cfg.Headers["Authorization"] = "Bearer " + cfg.BearerToken
    } else if cfg.BasicAuth != "" {
        parts := strings.SplitN(cfg.BasicAuth, ":", 2)
        if len(parts) != 2 {
            return fmt.Errorf("basic auth format should be 'user:pass'")
        }
        auth := base64.StdEncoding.EncodeToString([]byte(cfg.BasicAuth))
        cfg.Headers["Authorization"] = "Basic " + auth
    }
    return nil
}

// ExecuteRequest executes an HTTP request based on the provided config
func ExecuteRequest(cfg *config.Config) (*Response, error) {
    if err := addAuth(cfg); err != nil {
        return nil, err
    }

    httpClient := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
    if !cfg.FollowRedirects {
        httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        }
    }

    var body io.Reader
    if cfg.Body != "" {
        body = bytes.NewBufferString(cfg.Body)
    }

    req, err := http.NewRequest(cfg.Method, cfg.URL, body)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    for k, v := range cfg.Headers {
        req.Header.Set(k, v)
    }

    if cfg.Body != "" && req.Header.Get("Content-Type") == "" {
        req.Header.Set("Content-Type", "application/json")
    }

    start := time.Now()
    resp, err := httpClient.Do(req)
    duration := time.Since(start)

    if err != nil {
        return nil, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

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

