package cmd

import (
	"os"
	"strings"
	"testing"

	"logmd/vault"
)

// TestIsValidDateFormat tests the date format validation function.
// Learn: Testing edge cases and boundary conditions is crucial for robust software.
func TestIsValidDateFormat(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "ValidDate",
			input:    "2024-01-15",
			expected: true,
		},
		{
			name:     "ValidLeapYear",
			input:    "2024-02-29",
			expected: true,
		},
		{
			name:     "InvalidLeapYear",
			input:    "2023-02-29",
			expected: false,
		},
		{
			name:     "InvalidMonth",
			input:    "2024-13-01",
			expected: false,
		},
		{
			name:     "InvalidDay",
			input:    "2024-01-32",
			expected: false,
		},
		{
			name:     "WrongFormat",
			input:    "24-01-15",
			expected: false,
		},
		{
			name:     "SlashesInsteadOfDashes",
			input:    "2024/01/15",
			expected: false,
		},
		{
			name:     "NoLeadingZero",
			input:    "2024-1-5",
			expected: false,
		},
		{
			name:     "EmptyString",
			input:    "",
			expected: false,
		},
		{
			name:     "TooLong",
			input:    "2024-01-15-extra",
			expected: false,
		},
		{
			name:     "ContainsLetters",
			input:    "2024-ab-cd",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidDateFormat(tc.input)
			if result != tc.expected {
				t.Errorf("isValidDateFormat(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

// TestRunViewCommand tests the view command with valid entries.
func TestRunViewCommand(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "logmd-view-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test vault and entry
	v, err := vault.New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}

	testDate := "2024-01-15"
	testContent := `# Test Entry

This is a **test** entry with:

- Bullet points
- *Italic text*
- ~~Strikethrough~~

## Code Example

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `

> This is a blockquote with important information.

| Column 1 | Column 2 |
|----------|----------|
| Data 1   | Data 2   |

That's all for today!`

	err = v.WriteEntry(testDate, []byte(testContent))
	if err != nil {
		t.Fatalf("Failed to write test entry: %v", err)
	}

	// Save original environment
	originalDir := os.Getenv("LOGMD_DIRECTORY")
	defer func() {
		if originalDir != "" {
			os.Setenv("LOGMD_DIRECTORY", originalDir)
		} else {
			os.Unsetenv("LOGMD_DIRECTORY")
		}
	}()

	// Set test environment
	os.Setenv("LOGMD_DIRECTORY", tmpDir)

	// Test the view command
	err = runViewCommand(nil, []string{testDate})
	if err != nil {
		t.Fatalf("runViewCommand() failed: %v", err)
	}

	// Note: We can't easily capture the output without modifying the function,
	// but we can test that it completes without error
}

// TestRunViewCommandWithNonexistentEntry tests error handling for missing entries.
func TestRunViewCommandWithNonexistentEntry(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "logmd-view-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save original environment
	originalDir := os.Getenv("LOGMD_DIRECTORY")
	defer func() {
		if originalDir != "" {
			os.Setenv("LOGMD_DIRECTORY", originalDir)
		} else {
			os.Unsetenv("LOGMD_DIRECTORY")
		}
	}()

	// Set test environment
	os.Setenv("LOGMD_DIRECTORY", tmpDir)

	// Test with nonexistent entry
	err = runViewCommand(nil, []string{"2024-01-15"})
	if err == nil {
		t.Error("Expected error for nonexistent entry, got nil")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected 'does not exist' in error message, got: %v", err)
	}
}

// TestRunViewCommandWithInvalidDate tests error handling for invalid date formats.
func TestRunViewCommandWithInvalidDate(t *testing.T) {
	invalidDates := []string{
		"2024-13-01", // Invalid month
		"2024-01-32", // Invalid day
		"24-01-15",   // Wrong year format
		"2024/01/15", // Wrong separator
		"2024-1-5",   // Missing leading zeros
		"not-a-date", // Completely invalid
		"",           // Empty string
	}

	for _, invalidDate := range invalidDates {
		t.Run("InvalidDate_"+invalidDate, func(t *testing.T) {
			err := runViewCommand(nil, []string{invalidDate})
			if err == nil {
				t.Errorf("Expected error for invalid date %q, got nil", invalidDate)
			}

			if !strings.Contains(err.Error(), "invalid date format") {
				t.Errorf("Expected 'invalid date format' in error message, got: %v", err)
			}
		})
	}
}

// TestRunViewCommandWithInvalidDirectory tests error handling with bad directory.
func TestRunViewCommandWithInvalidDirectory(t *testing.T) {
	// Save original environment
	originalDir := os.Getenv("LOGMD_DIRECTORY")
	defer func() {
		if originalDir != "" {
			os.Setenv("LOGMD_DIRECTORY", originalDir)
		} else {
			os.Unsetenv("LOGMD_DIRECTORY")
		}
	}()

	// Create a file (not directory) to test error handling
	tmpFile, err := os.CreateTemp("", "logmd-view-invalid-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Try to use the file as a directory
	os.Setenv("LOGMD_DIRECTORY", tmpFile.Name())

	err = runViewCommand(nil, []string{"2024-01-15"})
	if err == nil {
		t.Error("Expected error with invalid directory, got nil")
	}

	// Should fail during vault initialization
	if !strings.Contains(err.Error(), "failed to initialize journal directory") {
		t.Errorf("Expected vault initialization error, got: %v", err)
	}
}

// TestViewCommandHelp tests that the help text is properly formatted.
func TestViewCommandHelp(t *testing.T) {
	// Get the help text
	helpOutput := viewCmd.Long

	// Check that key information is present
	expectedContent := []string{
		"glamour",
		"YYYY-MM-DD",
		"Examples:",
		"logmd view",
		"Colored headings",
		"Syntax-highlighted",
		"tables and lists",
	}

	for _, content := range expectedContent {
		if !strings.Contains(helpOutput, content) {
			t.Errorf("Help text should contain %q", content)
		}
	}
}

// TestViewCommandRegistration tests that the command is properly registered.
func TestViewCommandRegistration(t *testing.T) {
	// Check that view command exists in root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "view" {
			found = true
			break
		}
	}

	if !found {
		t.Error("view command should be registered with root command")
	}

	// Check basic command properties
	if viewCmd.Use != "view <YYYY-MM-DD>" {
		t.Errorf("Expected Use to be 'view <YYYY-MM-DD>', got %q", viewCmd.Use)
	}

	if viewCmd.Short == "" {
		t.Error("view command should have a short description")
	}

	if viewCmd.Long == "" {
		t.Error("view command should have a long description")
	}
}

// TestViewCommandArgs tests argument validation.
func TestViewCommandArgs(t *testing.T) {
	// Test that command requires exactly one argument
	err := viewCmd.Args(viewCmd, []string{})
	if err == nil {
		t.Error("Expected error with no arguments")
	}

	err = viewCmd.Args(viewCmd, []string{"2024-01-15", "extra"})
	if err == nil {
		t.Error("Expected error with too many arguments")
	}

	err = viewCmd.Args(viewCmd, []string{"2024-01-15"})
	if err != nil {
		t.Errorf("Expected no error with one argument, got: %v", err)
	}
}
