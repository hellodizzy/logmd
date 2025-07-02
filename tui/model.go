// Package tui provides the Bubble Tea terminal user interface for logmd timeline.
// This package implements an interactive browser for journal entries using
// the Elm architecture pattern with models, views, and updates.
//
// Learn: The Bubble Tea framework follows the Elm architecture pattern.
// See: https://github.com/charmbracelet/bubbletea#the-elm-architecture
package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"logmd/vault"
)

// Entry represents a single journal entry in the timeline.
// Learn: Struct fields should be exported (capitalized) when accessed outside the package.
// See: https://go.dev/tour/basics/3
type Entry struct {
	// Date is the entry date in YYYY-MM-DD format
	Date string
	// Path is the absolute file path to the entry
	Path string
	// Title is extracted from the first heading or "(untitled)"
	Title string
	// Preview contains the first few lines for expanded view
	Preview []string
	// Expanded indicates whether this entry is currently expanded
	Expanded bool
}

// Model holds the state for the timeline TUI.
// Learn: Bubble Tea models contain all the state needed for the interface.
// See: https://github.com/charmbracelet/bubbletea/blob/master/examples/simple/main.go
type Model struct {
	// entries contains all journal entries loaded from the vault
	entries []Entry
	// cursor tracks the currently selected entry index
	cursor int
	// viewport height for scrolling calculations
	viewportHeight int
	// scrollOffset for handling long lists
	scrollOffset int
	// quitting indicates the user wants to exit
	quitting bool
	// loading indicates entries are being loaded
	loading bool
	// err holds any error that occurred during operation
	err error
	// vaultDir is the directory containing journal entries
	vaultDir string
	// previewLines is the number of lines to show in previews
	previewLines int
}

// KeyMap defines keybindings for the timeline interface.
// Learn: Key maps in Bubble Tea provide consistent keyboard shortcuts.
// See: https://github.com/charmbracelet/bubbles/tree/master/key
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Toggle   key.Binding
	Quit     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
}

// DefaultKeyMap returns the default keybindings for timeline navigation.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "toggle expand"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdown", "page down"),
		),
	}
}

// NewModel creates a new timeline model with the specified vault directory and preview lines.
// Learn: Constructor functions should accept necessary configuration parameters.
func NewModel(vaultDir string, previewLines int) Model {
	return Model{
		entries:        []Entry{},
		cursor:         0,
		viewportHeight: 20, // Default height, will be updated on resize
		scrollOffset:   0,
		quitting:       false,
		loading:        true,
		err:            nil,
		vaultDir:       vaultDir,
		previewLines:   previewLines,
	}
}

// Error returns any error that occurred during operation.
// Learn: Error methods allow callers to check for errors after operations complete.
func (m Model) Error() error {
	return m.err
}

// LoadEntriesMsg is sent when entries have been loaded from the vault.
// Learn: Messages in Bubble Tea carry data between updates and commands.
// See: https://github.com/charmbracelet/bubbletea#commands
type LoadEntriesMsg struct {
	Entries []Entry
	Error   error
}

// LoadEntriesCmd returns a command that loads entries from the vault.
// This is called asynchronously to avoid blocking the UI.
func LoadEntriesCmd(vaultDir string, previewLines int) tea.Cmd {
	return func() tea.Msg {
		entries, err := loadEntriesFromVault(vaultDir, previewLines)
		return LoadEntriesMsg{
			Entries: entries,
			Error:   err,
		}
	}
}

// loadEntriesFromVault loads all journal entries from the vault directory.
// Learn: Helper functions should handle complex operations to keep main logic clean.
func loadEntriesFromVault(vaultDir string, previewLines int) ([]Entry, error) {
	// Create vault instance
	v, err := vault.New(vaultDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize vault: %w", err)
	}

	// Get list of all entries (returns filenames like "2024-01-01.md")
	entryFiles, err := v.ListEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}

	// Convert to Entry structs with previews
	entries := make([]Entry, 0, len(entryFiles))
	for _, filename := range entryFiles {
		// Strip .md extension to get date
		date := strings.TrimSuffix(filename, ".md")
		entry, err := createEntryFromDate(v, date, previewLines)
		if err != nil {
			// Log error but continue with other entries
			fmt.Fprintf(os.Stderr, "Warning: failed to load entry %s: %v\n", date, err)
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// createEntryFromDate creates an Entry struct from a date by reading the file.
// Learn: Small helper functions make code more readable and testable.
func createEntryFromDate(v *vault.Vault, date string, previewLines int) (Entry, error) {
	// Read entry content
	content, err := v.ReadEntry(date)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to read entry: %w", err)
	}

	// Extract title and preview
	title, preview := extractTitleAndPreview(string(content), previewLines)

	// Get file path
	entryPath := v.DatePath(date)

	return Entry{
		Date:     date,
		Path:     entryPath,
		Title:    title,
		Preview:  preview,
		Expanded: false,
	}, nil
}

// extractTitleAndPreview extracts the title and preview lines from entry content.
// Learn: Text processing functions are common in CLI applications.
func extractTitleAndPreview(content string, previewLines int) (string, []string) {
	lines := strings.Split(content, "\n")

	title := "(untitled)"
	var preview []string
	previewStart := 0

	// Extract title from first heading
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			title = strings.TrimSpace(trimmed[2:])
			previewStart = i + 1
			break
		}
	}

	// Extract preview lines (skip empty lines at start)
	previewCount := 0
	for i := previewStart; i < len(lines) && previewCount < previewLines; i++ {
		line := lines[i]
		if strings.TrimSpace(line) != "" || previewCount > 0 {
			preview = append(preview, line)
			previewCount++
		}
	}

	return title, preview
}

// Init returns the initial command for the model.
// Learn: Init is called once when the program starts.
func (m Model) Init() tea.Cmd {
	return LoadEntriesCmd(m.vaultDir, m.previewLines)
}
