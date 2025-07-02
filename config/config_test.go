package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoad verifies that configuration loading works with defaults.
// Learn: Configuration testing often involves temporary files and environment variables.
// See: https://pkg.go.dev/os#Setenv
func TestLoad(t *testing.T) {
	// Save original environment
	originalEditor := os.Getenv("EDITOR")
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
	}()

	// Test with default values
	config, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if config.Directory == "" {
		t.Error("Directory should not be empty")
	}

	if config.Editor == "" {
		t.Error("Editor should not be empty")
	}

	if config.PreviewLines <= 0 {
		t.Error("PreviewLines should be positive")
	}

	// Test default preview lines
	expectedPreviewLines := 5
	if config.PreviewLines != expectedPreviewLines {
		t.Errorf("Expected PreviewLines=%d, got %d", expectedPreviewLines, config.PreviewLines)
	}
}

// TestLoadWithEnvironment verifies that environment variables override defaults.
func TestLoadWithEnvironment(t *testing.T) {
	// Save original environment
	originalEditor := os.Getenv("LOGMD_EDITOR")
	originalDir := os.Getenv("LOGMD_DIRECTORY")
	originalPreview := os.Getenv("LOGMD_PREVIEW_LINES")

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
		if originalPreview != "" {
			os.Setenv("LOGMD_PREVIEW_LINES", originalPreview)
		} else {
			os.Unsetenv("LOGMD_PREVIEW_LINES")
		}
	}()

	// Set test environment variables
	testEditor := "nano"
	testDir := "/tmp/test-journal"
	testPreview := "10"

	os.Setenv("LOGMD_EDITOR", testEditor)
	os.Setenv("LOGMD_DIRECTORY", testDir)
	os.Setenv("LOGMD_PREVIEW_LINES", testPreview)

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if config.Editor != testEditor {
		t.Errorf("Expected Editor=%s, got %s", testEditor, config.Editor)
	}

	if config.Directory != testDir {
		t.Errorf("Expected Directory=%s, got %s", testDir, config.Directory)
	}

	if config.PreviewLines != 10 {
		t.Errorf("Expected PreviewLines=10, got %d", config.PreviewLines)
	}
}

// TestGetDefaultEditor verifies editor selection logic.
func TestGetDefaultEditor(t *testing.T) {
	// Save original environment
	originalEditor := os.Getenv("EDITOR")
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
	}()

	// Test with EDITOR set
	testEditor := "emacs"
	os.Setenv("EDITOR", testEditor)

	editor := getDefaultEditor()
	if editor != testEditor {
		t.Errorf("Expected editor=%s, got %s", testEditor, editor)
	}

	// Test with EDITOR unset
	os.Unsetenv("EDITOR")

	editor = getDefaultEditor()
	expectedDefault := "vim"
	if editor != expectedDefault {
		t.Errorf("Expected default editor=%s, got %s", expectedDefault, editor)
	}
}

// TestGetConfigPath verifies config path resolution.
func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()

	// Should either return empty string or a valid path
	if path != "" {
		if !filepath.IsAbs(path) {
			t.Error("Config path should be absolute")
		}

		expectedFilename := ".logmdconfig"
		if filepath.Base(path) != expectedFilename {
			t.Errorf("Expected config filename %s, got %s", expectedFilename, filepath.Base(path))
		}
	}
}
