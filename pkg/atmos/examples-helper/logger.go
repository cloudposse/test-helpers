package examples_helper

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var statusStyles = map[string]lipgloss.Style{
	"skipped":   lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Italic(true),
	"completed": lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
	"started":   lipgloss.NewStyle().Foreground(lipgloss.Color("14")),
	"failed":    lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Underline(true),
}

// Fallback style for unknown statuses
var fallbackStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

// Function to log the phase and status with appropriate styling
func (s *TestSuite) logPhaseStatus(phaseName, status string) {
	// Get the style for the given status, fallback if not found
	style, ok := statusStyles[status]
	if !ok {
		style = fallbackStyle
	}

	// Create the styled output
	output := lipgloss.NewStyle().Bold(true).Underline(true).Render(phaseName) + " → " + style.Render(status)

	if status == "skipped" {
		log.WithPrefix(s.T().Name()).Warn(output)
	} else if status == "failed" {
		log.WithPrefix(s.T().Name()).Error(output)
	} else {
		log.WithPrefix(s.T().Name()).Info(output)
	}
}
