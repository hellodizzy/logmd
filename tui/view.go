package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles for the timeline interface
// Learn: lipgloss provides a CSS-like API for terminal styling in Go.
// See: https://github.com/charmbracelet/lipgloss#usage
var (
	// Base styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	dateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Bold(true)

	iconStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	previewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#374151")).
			Padding(0, 2).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Padding(1, 0)
)

// View renders the timeline interface.
// Learn: View functions in Bubble Tea return strings that represent the UI.
// See: https://github.com/charmbracelet/bubbletea#view
func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.loading {
		return "Loading journal entries..."
	}

	if len(m.entries) == 0 {
		return "No journal entries found. Use 'logmd today' to create your first entry."
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("📖 Journal Timeline"))
	b.WriteString("\n\n")

	// Entries
	start, end := m.visibleRange()
	for i := start; i <= end && i < len(m.entries); i++ {
		entry := m.entries[i]
		b.WriteString(m.renderEntry(entry, i == m.cursor))
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/k up • ↓/j down • enter/space toggle • q quit"))

	return b.String()
}

// renderEntry renders a single timeline entry.
// Learn: Helper methods should handle specific rendering concerns for clarity.
func (m Model) renderEntry(entry Entry, selected bool) string {
	var b strings.Builder

	// Icon and date
	icon := iconStyle.Render("📅")
	date := dateStyle.Render(entry.Date)
	title := entry.Title

	line := fmt.Sprintf("%s %s %s", icon, date, title)

	if selected {
		line = selectedStyle.Render(line)
	} else {
		line = lipgloss.NewStyle().Padding(0, 1).Render(line)
	}

	b.WriteString(line)

	// Preview if expanded
	if entry.Expanded && len(entry.Preview) > 0 {
		b.WriteString("\n")
		for _, previewLine := range entry.Preview {
			if strings.TrimSpace(previewLine) != "" {
				b.WriteString(previewStyle.Render("  " + previewLine))
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

// visibleRange calculates which entries should be visible given the current scroll.
// Learn: Viewport calculations are important for performance with large lists.
func (m Model) visibleRange() (start, end int) {
	start = m.scrollOffset
	end = start + m.viewportHeight - 4 // Account for title and help text

	if end >= len(m.entries) {
		end = len(m.entries) - 1
	}

	return start, end
}
