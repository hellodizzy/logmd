package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"logmd/vault"
)

// TestRunTodayCommand tests the core today command functionality.
// Learn: Integration tests verify that multiple components work together correctly.
// See: https://martinfowler.com/articles/practical-test-pyramid.html
func TestRunTodayCommand(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "logmd-today-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save original environment
	originalEditor := os.Getenv("LOGMD_EDITOR")
	originalDir := os.Getenv("LOGMD_DIRECTORY")
	defer func() {
		if originalEditor != "" {
			os.Setenv("LOGMD_EDITOR", originalEditor)
		} else {
			os.Unsetenv("LOGMD_EDITOR")
		}
		if originalDir != "" {
			os.Setenv("LOGMD_DIRECTORY", originalDir)
		} else {
			os.Unsetenv("LOGMD_DIRECTORY")
		}
	}()

	// Set test environment - use 'true' as editor (always succeeds, exits quickly)
	os.Setenv("LOGMD_EDITOR", "true")
	os.Setenv("LOGMD_DIRECTORY", tmpDir)

	// Test creating new entry
	t.Run("CreateNewEntry", func(t *testing.T) {
		err := runTodayCommand(nil, []string{})
		if err != nil {
			t.Fatalf("runTodayCommand() failed: %v", err)
		}

		// Verify entry was created
		v, err := vault.New(tmpDir)
		if err != nil {
			t.Fatalf("Failed to create vault: %v", err)
		}

		if !v.TodayExists() {
			t.Error("Today's entry should exist after running command")
		}

		// Verify content
		today := time.Now().Format("2006-01-02")
		content, err := v.ReadEntry(today)
		if err != nil {
			t.Fatalf("Failed to read created entry: %v", err)
		}

		expectedContent := "# " + today + "\n\n"
		if string(content) != expectedContent {
			t.Errorf("Expected content %q, got %q", expectedContent, string(content))
		}
	})

	// Test opening existing entry
	t.Run("OpenExistingEntry", func(t *testing.T) {
		// Entry should already exist from previous test
		err := runTodayCommand(nil, []string{})
		if err != nil {
			t.Fatalf("runTodayCommand() failed on existing entry: %v", err)
		}

		// Entry should still exist and be unchanged
		v, err := vault.New(tmpDir)
		if err != nil {
			t.Fatalf("Failed to create vault: %v", err)
		}

		if !v.TodayExists() {
			t.Error("Today's entry should still exist")
		}
	})
}

// TestRunTodayCommandWithInvalidEditor tests error handling with bad editor.
func TestRunTodayCommandWithInvalidEditor(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "logmd-today-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save original environment
	originalEditor := os.Getenv("LOGMD_EDITOR")
	originalDir := os.Getenv("LOGMD_DIRECTORY")
	defer func() {
		if originalEditor != "" {
			os.Setenv("LOGMD_EDITOR", originalEditor)
		} else {
			os.Unsetenv("LOGMD_EDITOR")
		}
		if originalDir != "" {
			os.Setenv("LOGMD_DIRECTORY", originalDir)
		} else {
			os.Unsetenv("LOGMD_DIRECTORY")
		}
	}()

	// Set test environment with invalid editor
	os.Setenv("LOGMD_EDITOR", "nonexistent-editor-command")
	os.Setenv("LOGMD_DIRECTORY", tmpDir)

	// Test should fail with appropriate error
	err = runTodayCommand(nil, []string{})
	if err == nil {
		t.Error("Expected error with invalid editor, got nil")
	}

	if !strings.Contains(err.Error(), "failed to launch editor") {
		t.Errorf("Expected 'failed to launch editor' in error, got: %v", err)
	}
}

// TestRunTodayCommandWithInvalidDirectory tests error handling with bad directory.
func TestRunTodayCommandWithInvalidDirectory(t *testing.T) {
	// Save original environment
	originalEditor := os.Getenv("LOGMD_EDITOR")
	originalDir := os.Getenv("LOGMD_DIRECTORY")
	defer func() {
		if originalEditor != "" {
			os.Setenv("LOGMD_EDITOR", originalEditor)
		} else {
			os.Unsetenv("LOGMD_EDITOR")
		}
		if originalDir != "" {
			os.Setenv("LOGMD_DIRECTORY", originalDir)
		} else {
			os.Unsetenv("LOGMD_DIRECTORY")
		}
	}()

	// Set test environment with invalid directory (try to create under a file instead of directory)
	tmpFile, err := os.CreateTemp("", "logmd-invalid-dir-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	invalidDir := filepath.Join(tmpFile.Name(), "subdir") // This should fail since tmpFile.Name() is a file, not directory

	os.Setenv("LOGMD_EDITOR", "true")
	os.Setenv("LOGMD_DIRECTORY", invalidDir)

	// Test should fail with appropriate error
	err = runTodayCommand(nil, []string{})
	if err == nil {
		t.Error("Expected error with invalid directory, got nil")
		return
	}

	if !strings.Contains(err.Error(), "failed to initialize journal directory") {
		t.Errorf("Expected 'failed to initialize journal directory' in error, got: %v", err)
	}
}

// TestLaunchEditor tests the editor launching functionality.
func TestLaunchEditor(t *testing.T) {
	// Create temporary file for testing
	tmpFile, err := os.CreateTemp("", "logmd-editor-test-*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	testCases := []struct {
		name        string
		editor      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "ValidEditor",
			editor:      "true", // 'true' command always succeeds
			expectError: false,
		},
		{
			name:        "EditorExitsNonZero",
			editor:      "false", // 'false' command always fails with exit code 1
			expectError: true,
			errorMsg:    "editor exited with status",
		},
		{
			name:        "NonexistentEditor",
			editor:      "nonexistent-editor-command-12345",
			expectError: true,
			errorMsg:    "failed to run editor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := launchEditor(tc.editor, tmpFile.Name())

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for editor %s, got nil", tc.editor)
				} else if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Expected error containing %q, got: %v", tc.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for editor %s, got: %v", tc.editor, err)
				}
			}
		})
	}
}

// TestTodayCommandIntegration tests the full command integration including config loading.
func TestTodayCommandIntegration(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "logmd-today-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary config file
	configDir, err := os.MkdirTemp("", "logmd-config-*")
	if err != nil {
		t.Fatalf("Failed to create config temp dir: %v", err)
	}
	defer os.RemoveAll(configDir)

	configFile := filepath.Join(configDir, ".logmdconfig")
	configContent := `directory = "` + tmpDir + `"
editor = "true"
preview_lines = 5`

	err = os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Save original environment and HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// Set HOME to our config directory so the config file is found
	os.Setenv("HOME", configDir)

	// Test the command
	err = runTodayCommand(nil, []string{})
	if err != nil {
		t.Fatalf("runTodayCommand() failed with config file: %v", err)
	}

	// Verify entry was created in the configured directory
	v, err := vault.New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}

	if !v.TodayExists() {
		t.Error("Today's entry should exist after running command with config")
	}

	// Verify the entry has correct content
	today := time.Now().Format("2006-01-02")
	content, err := v.ReadEntry(today)
	if err != nil {
		t.Fatalf("Failed to read entry: %v", err)
	}

	expectedContent := "# " + today + "\n\n"
	if string(content) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(content))
	}
}
