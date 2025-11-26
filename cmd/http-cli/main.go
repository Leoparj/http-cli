package main

import (
	"fmt"
	"os"

	"github.com/I-invincib1e/http-cli/internal/client"
	"github.com/I-invincib1e/http-cli/internal/config"
	"github.com/I-invincib1e/http-cli/internal/output"
	"github.com/I-invincib1e/http-cli/internal/styles"
)

func main() {
	// Parse command-line flags
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// Initialize styles
	st := styles.New()

	// Display request information
	output.DisplayRequest(cfg, st)

	// Execute HTTP request
	resp, err := client.ExecuteRequest(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Display response
	if err := output.DisplayResponse(cfg, resp, st); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

