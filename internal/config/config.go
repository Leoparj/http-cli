package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the HTTP request
type Config struct {
	Method          string
	URL             string
	Headers         map[string]string
	Body            string
	BodyFile        string
	Timeout         int
	FollowRedirects bool
	BearerToken     string
	BasicAuth       string
	OutputFile      string
	Quiet           bool
	Verbose         bool
	StatusOnly      bool
}

// ParseFlags parses command-line flags and returns a Config struct
func ParseFlags() (*Config, error) {
	var method, url, body, headers, bodyFile, bearerToken, basicAuth, outputFile string
	var timeout int
	var followRedirects, quiet, verbose, statusOnly bool

	flag.StringVar(&method, "m", "GET", "HTTP method (GET, POST, PUT, DELETE, PATCH)")
	flag.StringVar(&method, "method", "GET", "HTTP method (GET, POST, PUT, DELETE, PATCH)")
	flag.StringVar(&url, "u", "", "URL to request")
	flag.StringVar(&url, "url", "", "URL to request")
	flag.StringVar(&body, "d", "", "Request body (JSON string)")
	flag.StringVar(&body, "data", "", "Request body (JSON string)")
	flag.StringVar(&bodyFile, "f", "", "Read request body from file")
	flag.StringVar(&bodyFile, "file", "", "Read request body from file")
	flag.StringVar(&headers, "H", "", "Headers (format: 'Key:Value,Key2:Value2')")
	flag.StringVar(&headers, "header", "", "Headers (format: 'Key:Value,Key2:Value2')")
	flag.StringVar(&bearerToken, "b", "", "Bearer token (sets Authorization header)")
	flag.StringVar(&bearerToken, "bearer", "", "Bearer token (sets Authorization header)")
	flag.StringVar(&basicAuth, "a", "", "Basic auth (format: 'user:pass')")
	flag.StringVar(&basicAuth, "auth", "", "Basic auth (format: 'user:pass')")
	flag.StringVar(&outputFile, "o", "", "Save response body to file")
	flag.StringVar(&outputFile, "output", "", "Save response body to file")
	flag.IntVar(&timeout, "t", 30, "Request timeout in seconds")
	flag.IntVar(&timeout, "timeout", 30, "Request timeout in seconds")
	flag.BoolVar(&followRedirects, "L", false, "Follow redirects")
	flag.BoolVar(&followRedirects, "follow", false, "Follow redirects")
	flag.BoolVar(&quiet, "q", false, "Quiet mode (only show response body)")
	flag.BoolVar(&quiet, "quiet", false, "Quiet mode (only show response body)")
	flag.BoolVar(&verbose, "v", false, "Verbose mode (show more details)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode (show more details)")
	flag.BoolVar(&statusOnly, "s", false, "Show only status code")
	flag.BoolVar(&statusOnly, "status-only", false, "Show only status code")

	flag.Usage = PrintUsage

	flag.Parse()

	// Read body from file if specified
	finalBody := body
	if bodyFile != "" {
		fileContent, err := os.ReadFile(bodyFile)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		finalBody = string(fileContent)
	}

	config := &Config{
		Method:          strings.ToUpper(method),
		URL:             url,
		Headers:         ParseHeaders(headers),
		Body:            finalBody,
		BodyFile:        bodyFile,
		Timeout:         timeout,
		FollowRedirects: followRedirects,
		BearerToken:     bearerToken,
		BasicAuth:       basicAuth,
		OutputFile:      outputFile,
		Quiet:           quiet,
		Verbose:         verbose,
		StatusOnly:      statusOnly,
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("URL is required")
	}
	return nil
}

// ParseHeaders parses a comma-separated header string into a map
func ParseHeaders(headersStr string) map[string]string {
	headers := make(map[string]string)
	if headersStr == "" {
		return headers
	}

	pairs := strings.Split(headersStr, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}
	return headers
}

// PrintUsage prints the usage information
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "HTTP CLI - A colorful command-line HTTP client\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s -m GET -u https://api.github.com/users/octocat\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -m POST -u https://api.example.com/users -f body.json\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -u https://api.example.com/data -b token123\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -u https://api.example.com/data -a 'user:pass'\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -u https://api.example.com/data -o response.json\n", os.Args[0])
}

