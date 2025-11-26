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

    // Helper to register short and long flag names for string values
    registerStr := func(short, long string, target *string, defVal, usage string) {
        flag.StringVar(target, short, defVal, usage)
        flag.StringVar(target, long, defVal, usage)
    }
    // Helper to register short and long flag names for bool values
    registerBool := func(short, long string, target *bool, usage string) {
        flag.BoolVar(target, short, false, usage)
        flag.BoolVar(target, long, false, usage)
    }
    // Helper to register short and long flag names for int values
    registerInt := func(short, long string, target *int, defVal int, usage string) {
        flag.IntVar(target, short, defVal, usage)
        flag.IntVar(target, long, defVal, usage)
    }

    registerStr("m", "method", &method, "GET", "HTTP method (GET, POST, PUT, DELETE, PATCH)")
    registerStr("u", "url", &url, "", "URL to request")
    registerStr("d", "data", &body, "", "Request body (JSON string)")
    registerStr("f", "file", &bodyFile, "", "Read request body from file")
    registerStr("H", "header", &headers, "", "Headers (format: 'Key:Value,Key2:Value2')")
    registerStr("b", "bearer", &bearerToken, "", "Bearer token (sets Authorization header)")
    registerStr("a", "auth", &basicAuth, "", "Basic auth (format: 'user:pass')")
    registerStr("o", "output", &outputFile, "", "Save response body to file")
    registerInt("t", "timeout", &timeout, 30, "Request timeout in seconds")
    registerBool("L", "follow", &followRedirects, "Follow redirects")
    registerBool("q", "quiet", &quiet, "Quiet mode (only show response body)")
    registerBool("v", "verbose", &verbose, "Verbose mode (show more details)")
    registerBool("s", "status-only", &statusOnly, "Show only status code")

    flag.Usage = PrintUsage
    flag.Parse()

    finalBody := body
    if bodyFile != "" {
        fileContent, err := os.ReadFile(bodyFile)
        if err != nil {
            return nil, fmt.Errorf("error reading file: %w", err)
        }
        finalBody = string(fileContent)
    }

    return &Config{
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
    }, nil
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

