package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Color styles
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1)

	methodStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33"))

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	bodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)
)

type Config struct {
	Method        string
	URL           string
	Headers       map[string]string
	Body          string
	BodyFile      string
	Timeout       int
	FollowRedirects bool
	BearerToken   string
	BasicAuth     string
	OutputFile    string
	Quiet         bool
	Verbose       bool
	StatusOnly    bool
}

func main() {
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

	flag.Usage = func() {
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

	flag.Parse()

	if url == "" {
		fmt.Fprintf(os.Stderr, "%s\n", errorStyle.Render("Error: URL is required"))
		flag.Usage()
		os.Exit(1)
	}

	// Read body from file if specified
	finalBody := body
	if bodyFile != "" {
		fileContent, err := os.ReadFile(bodyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s\n", errorStyle.Render("Error reading file:"), err.Error())
			os.Exit(1)
		}
		finalBody = string(fileContent)
	}

	config := Config{
		Method:          strings.ToUpper(method),
		URL:             url,
		Headers:         parseHeaders(headers),
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

	makeRequest(config)
}

func parseHeaders(headersStr string) map[string]string {
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

func makeRequest(config Config) {
	// Add bearer token to headers if provided
	if config.BearerToken != "" {
		config.Headers["Authorization"] = "Bearer " + config.BearerToken
	}

	// Add basic auth if provided
	if config.BasicAuth != "" {
		parts := strings.SplitN(config.BasicAuth, ":", 2)
		if len(parts) == 2 {
			auth := base64.StdEncoding.EncodeToString([]byte(config.BasicAuth))
			config.Headers["Authorization"] = "Basic " + auth
		} else {
			fmt.Fprintf(os.Stderr, "%s %s\n", errorStyle.Render("Error:"), "Basic auth format should be 'user:pass'")
			os.Exit(1)
		}
	}

	// Display request info (skip if quiet mode)
	if !config.Quiet && !config.StatusOnly {
		fmt.Println()
		fmt.Println(headerStyle.Render("REQUEST"))
		fmt.Println()
		fmt.Printf("%s %s\n", methodStyle.Render(config.Method), urlStyle.Render(config.URL))
		fmt.Println()

		if len(config.Headers) > 0 {
			fmt.Println(keyStyle.Render("Headers:"))
			for k, v := range config.Headers {
				// Mask sensitive headers in non-verbose mode
				if !config.Verbose && (k == "Authorization" || strings.ToLower(k) == "authorization") {
					if strings.HasPrefix(v, "Bearer ") {
						fmt.Printf("  %s: %s\n", keyStyle.Render(k), valueStyle.Render("Bearer ***"))
					} else if strings.HasPrefix(v, "Basic ") {
						fmt.Printf("  %s: %s\n", keyStyle.Render(k), valueStyle.Render("Basic ***"))
					} else {
						fmt.Printf("  %s: %s\n", keyStyle.Render(k), valueStyle.Render(v))
					}
				} else {
					fmt.Printf("  %s: %s\n", keyStyle.Render(k), valueStyle.Render(v))
				}
			}
			fmt.Println()
		}

		if config.Body != "" && config.Verbose {
			fmt.Println(keyStyle.Render("Body:"))
			prettyBody := formatJSON(config.Body)
			fmt.Println(bodyStyle.Render(prettyBody))
			fmt.Println()
		}
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	// Configure redirect following
	if !config.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Create request
	var req *http.Request
	var err error

	if config.Body != "" {
		req, err = http.NewRequest(config.Method, config.URL, bytes.NewBufferString(config.Body))
	} else {
		req, err = http.NewRequest(config.Method, config.URL, nil)
	}

	if err != nil {
		fmt.Printf("%s %s\n", errorStyle.Render("Error:"), err.Error())
		os.Exit(1)
	}

	// Set headers
	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	// If body exists and no Content-Type header, default to application/json
	if config.Body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make request
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("%s %s\n", errorStyle.Render("Error:"), err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s %s\n", errorStyle.Render("Error:"), err.Error())
		os.Exit(1)
	}

	// Save to file if output file specified
	if config.OutputFile != "" {
		err := os.WriteFile(config.OutputFile, bodyBytes, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s\n", errorStyle.Render("Error saving file:"), err.Error())
			os.Exit(1)
		}
		if !config.Quiet {
			fmt.Printf("%s %s\n", successStyle.Render("Response saved to:"), config.OutputFile)
		}
	}

	// Status only mode
	if config.StatusOnly {
		statusColor := "196" // red
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			statusColor = "46" // green
		} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			statusColor = "226" // yellow
		}
		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Bold(true)
		fmt.Printf("%s\n", statusStyle.Render(fmt.Sprintf("%d", resp.StatusCode)))
		os.Exit(0)
	}

	// Quiet mode - only show body
	if config.Quiet {
		if len(bodyBytes) > 0 {
			bodyStr := string(bodyBytes)
			prettyBody := formatJSON(bodyStr)
			fmt.Println(prettyBody)
		}
		os.Exit(0)
	}

	// Display response
	fmt.Println(headerStyle.Render("RESPONSE"))
	fmt.Println()

	// Status
	statusColor := "196" // red
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		statusColor = "46" // green
	} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		statusColor = "226" // yellow
	}

	statusText := fmt.Sprintf("%d %s", resp.StatusCode, resp.Status)
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Bold(true)
	fmt.Printf("%s %s\n", statusStyle.Render("Status:"), statusStyle.Render(statusText))
	fmt.Printf("%s %s\n", keyStyle.Render("Time:"), valueStyle.Render(duration.String()))
	
	if config.Verbose {
		fmt.Printf("%s %s\n", keyStyle.Render("Size:"), valueStyle.Render(fmt.Sprintf("%d bytes", len(bodyBytes))))
	}
	fmt.Println()

	// Headers
	if len(resp.Header) > 0 && config.Verbose {
		fmt.Println(keyStyle.Render("Headers:"))
		for k, v := range resp.Header {
			fmt.Printf("  %s: %s\n", keyStyle.Render(k), valueStyle.Render(strings.Join(v, ", ")))
		}
		fmt.Println()
	}

	// Body
	if len(bodyBytes) > 0 {
		if !config.Verbose {
			fmt.Println(keyStyle.Render("Body:"))
		} else {
			fmt.Println(keyStyle.Render("Body:"))
		}
		bodyStr := string(bodyBytes)
		prettyBody := formatJSON(bodyStr)
		fmt.Println(bodyStyle.Render(prettyBody))
	}
	fmt.Println()
}

func formatJSON(jsonStr string) string {
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return jsonStr // Return as-is if not valid JSON
	}

	prettyJSON, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return jsonStr
	}

	return string(prettyJSON)
}

