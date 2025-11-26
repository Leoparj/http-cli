package styles

import "github.com/charmbracelet/lipgloss"

// Styles holds all the color styles for terminal output
type Styles struct {
	Header  lipgloss.Style
	Method  lipgloss.Style
	Status  lipgloss.Style
	Success lipgloss.Style
	Error   lipgloss.Style
	URL     lipgloss.Style
	Key     lipgloss.Style
	Value   lipgloss.Style
	Body    lipgloss.Style
}

// New creates and returns a new Styles instance with all styles initialized
func New() *Styles {
	return &Styles{
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1),
		Method: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")),
		URL: lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")),
		Key: lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")),
		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		Body: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),
	}
}

// StatusColor returns the appropriate color code for a given HTTP status code
func StatusColor(statusCode int) string {
	if statusCode >= 200 && statusCode < 300 {
		return "46" // green
	} else if statusCode >= 300 && statusCode < 400 {
		return "226" // yellow
	}
	return "196" // red
}

