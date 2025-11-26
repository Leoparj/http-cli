package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/I-invincib1e/http-cli/internal/client"
	"github.com/I-invincib1e/http-cli/internal/config"
	"github.com/I-invincib1e/http-cli/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

// DisplayRequest displays the request information
func DisplayRequest(cfg *config.Config, st *styles.Styles) {
	if cfg.Quiet || cfg.StatusOnly {
		return
	}

	fmt.Println()
	fmt.Println(st.Header.Render("REQUEST"))
	fmt.Println()
	fmt.Printf("%s %s\n", st.Method.Render(cfg.Method), st.URL.Render(cfg.URL))
	fmt.Println()

	if len(cfg.Headers) > 0 {
		fmt.Println(st.Key.Render("Headers:"))
		for k, v := range cfg.Headers {
			// Mask sensitive headers in non-verbose mode
			if !cfg.Verbose && (k == "Authorization" || strings.ToLower(k) == "authorization") {
				if strings.HasPrefix(v, "Bearer ") {
					fmt.Printf("  %s: %s\n", st.Key.Render(k), st.Value.Render("Bearer ***"))
				} else if strings.HasPrefix(v, "Basic ") {
					fmt.Printf("  %s: %s\n", st.Key.Render(k), st.Value.Render("Basic ***"))
				} else {
					fmt.Printf("  %s: %s\n", st.Key.Render(k), st.Value.Render(v))
				}
			} else {
				fmt.Printf("  %s: %s\n", st.Key.Render(k), st.Value.Render(v))
			}
		}
		fmt.Println()
	}

	if cfg.Body != "" && cfg.Verbose {
		fmt.Println(st.Key.Render("Body:"))
		prettyBody := FormatJSON(cfg.Body)
		fmt.Println(st.Body.Render(prettyBody))
		fmt.Println()
	}
}

// DisplayResponse displays the response information
func DisplayResponse(cfg *config.Config, resp *client.Response, st *styles.Styles) error {
	// Save to file if output file specified
	if cfg.OutputFile != "" {
		err := os.WriteFile(cfg.OutputFile, resp.Body, 0644)
		if err != nil {
			return fmt.Errorf("error saving file: %w", err)
		}
		if !cfg.Quiet {
			fmt.Printf("%s %s\n", st.Success.Render("Response saved to:"), cfg.OutputFile)
		}
	}

	// Status only mode
	if cfg.StatusOnly {
		statusColor := styles.StatusColor(resp.StatusCode)
		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Bold(true)
		fmt.Printf("%s\n", statusStyle.Render(fmt.Sprintf("%d", resp.StatusCode)))
		return nil
	}

	// Quiet mode - only show body
	if cfg.Quiet {
		if len(resp.Body) > 0 {
			bodyStr := string(resp.Body)
			prettyBody := FormatJSON(bodyStr)
			fmt.Println(prettyBody)
		}
		return nil
	}

	// Display response
	fmt.Println(st.Header.Render("RESPONSE"))
	fmt.Println()

	// Status
	statusColor := styles.StatusColor(resp.StatusCode)
	statusText := fmt.Sprintf("%d %s", resp.StatusCode, resp.Status)
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Bold(true)
	fmt.Printf("%s %s\n", statusStyle.Render("Status:"), statusStyle.Render(statusText))
	fmt.Printf("%s %s\n", st.Key.Render("Time:"), st.Value.Render(resp.Duration.String()))

	if cfg.Verbose {
		fmt.Printf("%s %s\n", st.Key.Render("Size:"), st.Value.Render(fmt.Sprintf("%d bytes", len(resp.Body))))
	}
	fmt.Println()

	// Headers
	if len(resp.Headers) > 0 && cfg.Verbose {
		fmt.Println(st.Key.Render("Headers:"))
		for k, v := range resp.Headers {
			fmt.Printf("  %s: %s\n", st.Key.Render(k), st.Value.Render(strings.Join(v, ", ")))
		}
		fmt.Println()
	}

	// Body
	if len(resp.Body) > 0 {
		fmt.Println(st.Key.Render("Body:"))
		bodyStr := string(resp.Body)
		prettyBody := FormatJSON(bodyStr)
		fmt.Println(st.Body.Render(prettyBody))
	}
	fmt.Println()

	return nil
}

// FormatJSON attempts to format a JSON string with indentation
func FormatJSON(jsonStr string) string {
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

