package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all state changes in response to messages.
// Learn: Update functions in Bubble Tea handle state transitions and side effects.
// See: https://github.com/charmbracelet/bubbletea#update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.viewportHeight = msg.Height - 6 // Account for title, help, and padding
		return m, nil

	case LoadEntriesMsg:
		m.loading = false
		if msg.Error != nil {
			m.err = msg.Error
			return m, nil
		}
		m.entries = msg.Entries
		return m, nil

	default:
		return m, nil
	}
}

// handleKeyPress processes keyboard input and returns updated model and commands.
// Learn: Switch statements on type assertions are a common Go pattern.
// See: https://go.dev/tour/methods/16
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.entries) == 0 {
		// Only allow quit when no entries
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.adjustScroll()
		}

	case "down", "j":
		if m.cursor < len(m.entries)-1 {
			m.cursor++
			m.adjustScroll()
		}

	case "enter", " ":
		if m.cursor < len(m.entries) {
			m.entries[m.cursor].Expanded = !m.entries[m.cursor].Expanded
		}

	case "pgup":
		m.cursor -= 10
		if m.cursor < 0 {
			m.cursor = 0
		}
		m.adjustScroll()

	case "pgdown":
		m.cursor += 10
		if m.cursor >= len(m.entries) {
			m.cursor = len(m.entries) - 1
		}
		m.adjustScroll()

	case "home":
		m.cursor = 0
		m.adjustScroll()

	case "end":
		m.cursor = len(m.entries) - 1
		m.adjustScroll()
	}

	return m, nil
}

// adjustScroll ensures the cursor is visible within the viewport.
// Learn: Scrolling logic requires careful bounds checking and offset management.
func (m *Model) adjustScroll() {
	visibleHeight := m.viewportHeight - 4 // Account for title and help

	// Scroll up if cursor is above viewport
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}

	// Scroll down if cursor is below viewport
	if m.cursor >= m.scrollOffset+visibleHeight {
		m.scrollOffset = m.cursor - visibleHeight + 1
	}

	// Ensure scroll offset is within bounds
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}

	maxScroll := len(m.entries) - visibleHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
	}
}
