package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestNew verifies that New creates a vault with the correct directory structure.
// Learn: Test functions must start with "Test" and take *testing.T as parameter.
// See: https://pkg.go.dev/testing#T
func TestNew(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if vault.Directory == "" {
		t.Error("Directory should not be empty")
	}

	// Verify directory exists and has correct permissions
	info, err := os.Stat(vault.Directory)
	if err != nil {
		t.Fatalf("Directory does not exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("Path should be a directory")
	}

	// Check permissions (should be 0700)
	if info.Mode().Perm() != 0700 {
		t.Errorf("Expected permissions 0700, got %v", info.Mode().Perm())
	}
}

// TestNewWithInvalidPath tests error handling when path cannot be resolved.
func TestNewWithInvalidPath(t *testing.T) {
	// Test with a path that contains null bytes (invalid on most systems)
	invalidPath := "/tmp/test\x00invalid"
	_, err := New(invalidPath)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

// TestTodayPath verifies that TodayPath returns the correct format.
// Learn: Subtests allow organizing related test cases using t.Run().
// See: https://pkg.go.dev/testing#T.Run
func TestTodayPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	todayPath := vault.TodayPath()
	expectedSuffix := time.Now().Format("2006-01-02.md")

	if !filepath.IsAbs(todayPath) {
		t.Error("TodayPath should return absolute path")
	}

	if filepath.Base(todayPath) != expectedSuffix {
		t.Errorf("Expected filename %s, got %s", expectedSuffix, filepath.Base(todayPath))
	}
}

// TestDatePath verifies that DatePath returns correct paths for specific dates.
func TestDatePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	testCases := []struct {
		date     string
		expected string
	}{
		{"2024-01-15", "2024-01-15.md"},
		{"2023-12-31", "2023-12-31.md"},
		{"2024-02-29", "2024-02-29.md"}, // Leap year
	}

	for _, tc := range testCases {
		t.Run(tc.date, func(t *testing.T) {
			path := vault.DatePath(tc.date)
			if filepath.Base(path) != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, filepath.Base(path))
			}
		})
	}
}

// TestEntryExists verifies entry existence checking.
func TestEntryExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	testDate := "2024-01-15"

	// Should not exist initially
	if vault.EntryExists(testDate) {
		t.Error("Entry should not exist initially")
	}

	// Create the file
	testPath := vault.DatePath(testDate)
	if err := os.WriteFile(testPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Should exist now
	if !vault.EntryExists(testDate) {
		t.Error("Entry should exist after creation")
	}
}

// TestTodayExists verifies today's entry existence checking.
func TestTodayExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Should not exist initially
	if vault.TodayExists() {
		t.Error("Today's entry should not exist initially")
	}

	// Create today's file
	todayPath := vault.TodayPath()
	if err := os.WriteFile(todayPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create today's file: %v", err)
	}

	// Should exist now
	if !vault.TodayExists() {
		t.Error("Today's entry should exist after creation")
	}
}

// TestReadEntry verifies reading entry content.
func TestReadEntry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	testDate := "2024-01-15"
	testContent := "# Test Entry\n\nThis is test content."

	// Test reading non-existent entry
	_, err = vault.ReadEntry(testDate)
	if err == nil {
		t.Error("Expected error when reading non-existent entry")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected 'does not exist' error, got: %v", err)
	}

	// Create the entry
	testPath := vault.DatePath(testDate)
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test reading existing entry
	content, err := vault.ReadEntry(testDate)
	if err != nil {
		t.Fatalf("Failed to read entry: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(content))
	}
}

// TestWriteEntry verifies writing entry content.
func TestWriteEntry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	testDate := "2024-01-15"
	testContent := "# Test Entry\n\nThis is test content."

	// Write the entry
	err = vault.WriteEntry(testDate, []byte(testContent))
	if err != nil {
		t.Fatalf("Failed to write entry: %v", err)
	}

	// Verify the file was created and has correct content
	testPath := vault.DatePath(testDate)
	content, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(content))
	}

	// Test overwriting
	newContent := "# Updated Entry\n\nThis is updated content."
	err = vault.WriteEntry(testDate, []byte(newContent))
	if err != nil {
		t.Fatalf("Failed to overwrite entry: %v", err)
	}

	content, err = os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read overwritten file: %v", err)
	}

	if string(content) != newContent {
		t.Errorf("Expected updated content %q, got %q", newContent, string(content))
	}
}

// TestCreateEntry verifies entry creation with template.
func TestCreateEntry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	testDate := "2024-01-15"

	// Create the entry
	err = vault.CreateEntry(testDate)
	if err != nil {
		t.Fatalf("Failed to create entry: %v", err)
	}

	// Verify the file was created with correct template
	content, err := vault.ReadEntry(testDate)
	if err != nil {
		t.Fatalf("Failed to read created entry: %v", err)
	}

	expectedContent := "# 2024-01-15\n\n"
	if string(content) != expectedContent {
		t.Errorf("Expected template %q, got %q", expectedContent, string(content))
	}

	// Test creating entry that already exists
	err = vault.CreateEntry(testDate)
	if err == nil {
		t.Error("Expected error when creating existing entry")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}
}

// TestCreateTodayEntry verifies today's entry creation.
func TestCreateTodayEntry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Create today's entry
	err = vault.CreateTodayEntry()
	if err != nil {
		t.Fatalf("Failed to create today's entry: %v", err)
	}

	// Verify today's entry exists
	if !vault.TodayExists() {
		t.Error("Today's entry should exist after creation")
	}

	// Verify content
	today := time.Now().Format("2006-01-02")
	content, err := vault.ReadEntry(today)
	if err != nil {
		t.Fatalf("Failed to read today's entry: %v", err)
	}

	expectedContent := "# " + today + "\n\n"
	if string(content) != expectedContent {
		t.Errorf("Expected template %q, got %q", expectedContent, string(content))
	}

	// Test creating today's entry when it already exists
	err = vault.CreateTodayEntry()
	if err == nil {
		t.Error("Expected error when creating existing today's entry")
	}
}

// TestGetEntryInfo verifies entry metadata retrieval.
func TestGetEntryInfo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	testDate := "2024-01-15"

	// Test non-existent entry
	info := vault.GetEntryInfo(testDate)
	if info.Date != testDate {
		t.Errorf("Expected date %s, got %s", testDate, info.Date)
	}
	if info.Exists {
		t.Error("Entry should not exist")
	}
	if info.Size != 0 {
		t.Errorf("Expected size 0, got %d", info.Size)
	}

	// Create the entry
	testContent := "# Test Entry\n\nThis is test content."
	err = vault.WriteEntry(testDate, []byte(testContent))
	if err != nil {
		t.Fatalf("Failed to write entry: %v", err)
	}

	// Test existing entry
	info = vault.GetEntryInfo(testDate)
	if info.Date != testDate {
		t.Errorf("Expected date %s, got %s", testDate, info.Date)
	}
	if !info.Exists {
		t.Error("Entry should exist")
	}
	if info.Size != int64(len(testContent)) {
		t.Errorf("Expected size %d, got %d", len(testContent), info.Size)
	}
	if info.ModTime.IsZero() {
		t.Error("ModTime should not be zero for existing file")
	}

	expectedPath := vault.DatePath(testDate)
	if info.Path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, info.Path)
	}
}

// TestListEntries verifies that ListEntries correctly identifies and sorts markdown files.
func TestListEntries(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Create test files
	testFiles := []string{
		"2024-01-15.md",
		"2024-01-10.md",
		"2024-01-20.md",
		"not-a-date.md",  // Should be ignored
		"2024-01-15.txt", // Wrong extension, should be ignored
		"README.md",      // Not a date format, should be ignored
	}

	for _, filename := range testFiles {
		path := filepath.Join(vault.Directory, filename)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	entries, err := vault.ListEntries()
	if err != nil {
		t.Fatalf("ListEntries() failed: %v", err)
	}

	// Should return only valid date files, sorted newest first
	expected := []string{
		"2024-01-20.md",
		"2024-01-15.md",
		"2024-01-10.md",
	}

	if len(entries) != len(expected) {
		t.Errorf("Expected %d entries, got %d", len(expected), len(entries))
	}

	for i, entry := range entries {
		if i < len(expected) && entry != expected[i] {
			t.Errorf("Entry %d: expected %s, got %s", i, expected[i], entry)
		}
	}
}

// TestListEntriesInfo verifies metadata listing for all entries.
func TestListEntriesInfo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Create test entries
	testDates := []string{"2024-01-15", "2024-01-10", "2024-01-20"}
	for _, date := range testDates {
		err = vault.CreateEntry(date)
		if err != nil {
			t.Fatalf("Failed to create entry %s: %v", date, err)
		}
	}

	entries, err := vault.ListEntriesInfo()
	if err != nil {
		t.Fatalf("ListEntriesInfo() failed: %v", err)
	}

	// Should return entries sorted newest first
	expectedDates := []string{"2024-01-20", "2024-01-15", "2024-01-10"}

	if len(entries) != len(expectedDates) {
		t.Errorf("Expected %d entries, got %d", len(expectedDates), len(entries))
	}

	for i, entry := range entries {
		if i < len(expectedDates) {
			if entry.Date != expectedDates[i] {
				t.Errorf("Entry %d: expected date %s, got %s", i, expectedDates[i], entry.Date)
			}
			if !entry.Exists {
				t.Errorf("Entry %d (%s) should exist", i, entry.Date)
			}
			if entry.Size <= 0 {
				t.Errorf("Entry %d (%s) should have positive size", i, entry.Date)
			}
		}
	}
}

// TestIsValidDateFormat verifies the date format validation function.
func TestIsValidDateFormat(t *testing.T) {
	testCases := []struct {
		filename string
		valid    bool
	}{
		{"2024-01-15.md", true},
		{"2023-12-31.md", true},
		{"2024-02-29.md", true},  // Leap year
		{"2023-02-29.md", false}, // Not a leap year
		{"2024-13-01.md", false}, // Invalid month
		{"2024-01-32.md", false}, // Invalid day
		{"not-a-date.md", false},
		{"2024-01-15.txt", false}, // Wrong extension
		{"README.md", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			result := isValidDateFormat(tc.filename)
			if result != tc.valid {
				t.Errorf("isValidDateFormat(%s) = %v, expected %v", tc.filename, result, tc.valid)
			}
		})
	}
}
