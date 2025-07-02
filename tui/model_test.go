package tui

import (
	"os"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"logmd/vault"
)

// TestNewModel tests the model constructor.
// Learn: Constructor tests should verify initial state is correct.
func TestNewModel(t *testing.T) {
	vaultDir := "/test/vault"
	previewLines := 10

	model := NewModel(vaultDir, previewLines)

	// Verify initial state
	if model.vaultDir != vaultDir {
		t.Errorf("Expected vaultDir %s, got %s", vaultDir, model.vaultDir)
	}

	if model.previewLines != previewLines {
		t.Errorf("Expected previewLines %d, got %d", previewLines, model.previewLines)
	}

	if len(model.entries) != 0 {
		t.Errorf("Expected empty entries, got %d", len(model.entries))
	}

	if model.cursor != 0 {
		t.Errorf("Expected cursor 0, got %d", model.cursor)
	}

	if !model.loading {
		t.Error("Expected loading to be true initially")
	}

	if model.err != nil {
		t.Errorf("Expected no error initially, got %v", model.err)
	}
}

// TestExtractTitleAndPreview tests title and preview extraction from content.
func TestExtractTitleAndPreview(t *testing.T) {
	testCases := []struct {
		name            string
		content         string
		previewLines    int
		expectedTitle   string
		expectedPreview []string
	}{
		{
			name:            "SimpleHeading",
			content:         "# Daily Journal\n\nToday was a good day.\nI learned something new.",
			previewLines:    2,
			expectedTitle:   "Daily Journal",
			expectedPreview: []string{"Today was a good day.", "I learned something new."},
		},
		{
			name:            "NoHeading",
			content:         "Just some content\nwithout a heading\nhere.",
			previewLines:    2,
			expectedTitle:   "(untitled)",
			expectedPreview: []string{"Just some content", "without a heading"},
		},
		{
			name:            "EmptyContent",
			content:         "",
			previewLines:    2,
			expectedTitle:   "(untitled)",
			expectedPreview: []string{},
		},
		{
			name:            "HeadingOnly",
			content:         "# Just a Title",
			previewLines:    2,
			expectedTitle:   "Just a Title",
			expectedPreview: []string{},
		},
		{
			name:            "WithEmptyLines",
			content:         "# Title\n\n\nFirst line\n\nSecond line",
			previewLines:    3,
			expectedTitle:   "Title",
			expectedPreview: []string{"First line", "", "Second line"},
		},
		{
			name:            "LimitPreviewLines",
			content:         "# Title\n\nLine 1\nLine 2\nLine 3\nLine 4",
			previewLines:    2,
			expectedTitle:   "Title",
			expectedPreview: []string{"Line 1", "Line 2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			title, preview := extractTitleAndPreview(tc.content, tc.previewLines)

			if title != tc.expectedTitle {
				t.Errorf("Expected title %q, got %q", tc.expectedTitle, title)
			}

			// Handle empty slice comparison properly
			if len(preview) == 0 && len(tc.expectedPreview) == 0 {
				// Both empty, this is expected
			} else if !reflect.DeepEqual(preview, tc.expectedPreview) {
				t.Errorf("Expected preview %v, got %v", tc.expectedPreview, preview)
			}
		})
	}
}

// TestLoadEntriesFromVault tests loading entries from a vault.
func TestLoadEntriesFromVault(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "logmd-tui-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create vault and test entries
	v, err := vault.New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}

	testEntries := map[string]string{
		"2024-01-01": "# New Year Resolution\n\nI want to write more.",
		"2024-01-02": "# Day Two\n\nContinuing the journey.\nFeeling motivated.",
		"2024-01-03": "No heading today\n\nJust some thoughts.",
	}

	for date, content := range testEntries {
		err = v.WriteEntry(date, []byte(content))
		if err != nil {
			t.Fatalf("Failed to write test entry %s: %v", date, err)
		}
	}

	// Test loading entries
	entries, err := loadEntriesFromVault(tmpDir, 2)
	if err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Should have 3 entries (vault.ListEntries returns newest first)
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Check that entries are loaded correctly
	// Note: vault.ListEntries returns entries in reverse chronological order
	expectedDates := []string{"2024-01-03", "2024-01-02", "2024-01-01"}
	expectedTitles := []string{"(untitled)", "Day Two", "New Year Resolution"}

	for i, entry := range entries {
		if entry.Date != expectedDates[i] {
			t.Errorf("Entry %d: expected date %s, got %s", i, expectedDates[i], entry.Date)
		}

		if entry.Title != expectedTitles[i] {
			t.Errorf("Entry %d: expected title %q, got %q", i, expectedTitles[i], entry.Title)
		}

		if entry.Expanded {
			t.Errorf("Entry %d: should not be expanded initially", i)
		}

		if len(entry.Preview) > 2 {
			t.Errorf("Entry %d: preview should be limited to 2 lines, got %d", i, len(entry.Preview))
		}
	}
}

// TestLoadEntriesFromVaultError tests error handling when vault loading fails.
func TestLoadEntriesFromVaultError(t *testing.T) {
	// Try to load from non-existent directory
	entries, err := loadEntriesFromVault("/nonexistent/directory", 5)

	if err == nil {
		t.Error("Expected error when loading from non-existent directory")
	}

	if entries != nil {
		t.Error("Expected nil entries when error occurs")
	}
}

// TestModelUpdate tests the model update function with various messages.
func TestModelUpdate(t *testing.T) {
	model := NewModel("/test", 5)

	// Test window size message
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, cmd := model.Update(windowMsg)
	if cmd != nil {
		t.Error("WindowSizeMsg should not return a command")
	}

	m := updatedModel.(Model)
	expectedHeight := 24 - 6 // Height minus padding
	if m.viewportHeight != expectedHeight {
		t.Errorf("Expected viewportHeight %d, got %d", expectedHeight, m.viewportHeight)
	}

	// Test LoadEntriesMsg with success
	entries := []Entry{
		{Date: "2024-01-01", Title: "Test", Preview: []string{"Preview"}, Expanded: false},
	}
	loadMsg := LoadEntriesMsg{Entries: entries, Error: nil}

	updatedModel, cmd = model.Update(loadMsg)
	if cmd != nil {
		t.Error("LoadEntriesMsg should not return a command")
	}

	m = updatedModel.(Model)
	if m.loading {
		t.Error("Model should not be loading after LoadEntriesMsg")
	}

	if len(m.entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(m.entries))
	}

	// Test LoadEntriesMsg with error
	loadErrorMsg := LoadEntriesMsg{Entries: nil, Error: os.ErrNotExist}
	updatedModel, cmd = model.Update(loadErrorMsg)

	m = updatedModel.(Model)
	if m.err != os.ErrNotExist {
		t.Errorf("Expected error %v, got %v", os.ErrNotExist, m.err)
	}
}

// TestModelInit tests the model initialization.
func TestModelInit(t *testing.T) {
	model := NewModel("/test/vault", 5)

	cmd := model.Init()
	if cmd == nil {
		t.Error("Init should return a command to load entries")
	}

	// Execute the command to get the message
	msg := cmd()

	// Should return a LoadEntriesMsg
	if loadMsg, ok := msg.(LoadEntriesMsg); ok {
		// Should have an error since /test/vault doesn't exist
		if loadMsg.Error == nil {
			t.Error("Expected error when loading from non-existent vault")
		}
	} else {
		t.Errorf("Expected LoadEntriesMsg, got %T", msg)
	}
}

// TestModelError tests the Error method.
func TestModelError(t *testing.T) {
	model := NewModel("/test", 5)

	// Initially no error
	if model.Error() != nil {
		t.Errorf("Expected no error initially, got %v", model.Error())
	}

	// Set an error and test
	model.err = os.ErrNotExist
	if model.Error() != os.ErrNotExist {
		t.Errorf("Expected error %v, got %v", os.ErrNotExist, model.Error())
	}
}
