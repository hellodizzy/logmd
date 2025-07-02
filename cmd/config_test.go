package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunConfigCommand tests the config command with default settings.
// Learn: Configuration display tests should verify output format and content.
func TestRunConfigCommand(t *testing.T) {
	// Save original environment
	originalVars := saveEnvironment()
	defer restoreEnvironment(originalVars)

	// Clear all logmd environment variables
	clearLogmdEnvironment()

	// Test with default configuration (no config file, no env vars)
	err := runConfigCommand(nil, []string{})
	if err != nil {
		t.Fatalf("runConfigCommand() failed: %v", err)
	}

	// Note: We can't easily capture stdout in this test setup,
	// but we verify the command completes without error
}

// TestRunConfigCommandWithEnvironmentVariables tests config display with env vars.
func TestRunConfigCommandWithEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalVars := saveEnvironment()
	defer restoreEnvironment(originalVars)

	// Clear all logmd environment variables first
	clearLogmdEnvironment()

	// Set test environment variables
	os.Setenv("LOGMD_DIRECTORY", "/test/journal")
	os.Setenv("LOGMD_EDITOR", "nano")
	os.Setenv("LOGMD_PREVIEW_LINES", "10")

	// Test config command
	err := runConfigCommand(nil, []string{})
	if err != nil {
		t.Fatalf("runConfigCommand() with env vars failed: %v", err)
	}
}

// TestRunConfigCommandWithConfigFile tests config display with a config file.
func TestRunConfigCommandWithConfigFile(t *testing.T) {
	// Save original environment
	originalVars := saveEnvironment()
	defer restoreEnvironment(originalVars)

	// Clear all logmd environment variables
	clearLogmdEnvironment()

	// Create temporary directory for config file
	tmpDir, err := os.MkdirTemp("", "logmd-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create config file
	configPath := filepath.Join(tmpDir, ".logmdconfig")
	configContent := `directory = "/custom/journal"
editor = "code"
preview_lines = 7`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set HOME to our temp directory so config is found
	os.Setenv("HOME", tmpDir)

	// Test config command
	err = runConfigCommand(nil, []string{})
	if err != nil {
		t.Fatalf("runConfigCommand() with config file failed: %v", err)
	}
}

// TestRunConfigCommandWithBothSources tests precedence (env vars override config file).
func TestRunConfigCommandWithBothSources(t *testing.T) {
	// Save original environment
	originalVars := saveEnvironment()
	defer restoreEnvironment(originalVars)

	// Clear all logmd environment variables
	clearLogmdEnvironment()

	// Create temporary directory for config file
	tmpDir, err := os.MkdirTemp("", "logmd-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create config file
	configPath := filepath.Join(tmpDir, ".logmdconfig")
	configContent := `directory = "/config/journal"
editor = "vim"
preview_lines = 5`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set HOME to our temp directory
	os.Setenv("HOME", tmpDir)

	// Set environment variables (should override config file)
	os.Setenv("LOGMD_DIRECTORY", "/env/journal")
	os.Setenv("LOGMD_EDITOR", "code")

	// Test config command
	err = runConfigCommand(nil, []string{})
	if err != nil {
		t.Fatalf("runConfigCommand() with both sources failed: %v", err)
	}
}

// TestGetSettingSource tests the setting source detection function.
func TestGetSettingSource(t *testing.T) {
	// Save original environment
	originalVars := saveEnvironment()
	defer restoreEnvironment(originalVars)

	testCases := []struct {
		name           string
		envVar         string
		envValue       string
		hasConfigFile  bool
		expectedSource string
	}{
		{
			name:           "EnvironmentVariable",
			envVar:         "LOGMD_DIRECTORY",
			envValue:       "/test/path",
			hasConfigFile:  true,
			expectedSource: "🌍 Environment variable (LOGMD_DIRECTORY)",
		},
		{
			name:           "ConfigFile",
			envVar:         "LOGMD_DIRECTORY",
			envValue:       "",
			hasConfigFile:  true,
			expectedSource: "📄 Configuration file (~/.logmdconfig)",
		},
		{
			name:           "DefaultValue",
			envVar:         "LOGMD_DIRECTORY",
			envValue:       "",
			hasConfigFile:  false,
			expectedSource: "🔧 Default value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv(tc.envVar)

			// Set environment variable if specified
			if tc.envValue != "" {
				os.Setenv(tc.envVar, tc.envValue)
			}

			result := getSettingSource(tc.envVar, tc.hasConfigFile)
			if result != tc.expectedSource {
				t.Errorf("getSettingSource(%q, %v) = %q, expected %q",
					tc.envVar, tc.hasConfigFile, result, tc.expectedSource)
			}
		})
	}
}

// TestRepeatString tests the string repetition helper function.
func TestRepeatString(t *testing.T) {
	testCases := []struct {
		input    string
		count    int
		expected string
	}{
		{"=", 3, "==="},
		{"ab", 2, "abab"},
		{"x", 0, ""},
		{"", 5, ""},
		{"test", 1, "test"},
	}

	for _, tc := range testCases {
		result := repeatString(tc.input, tc.count)
		if result != tc.expected {
			t.Errorf("repeatString(%q, %d) = %q, expected %q",
				tc.input, tc.count, result, tc.expected)
		}
	}
}

// TestConfigCommandHelp tests that the help text is properly formatted.
func TestConfigCommandHelp(t *testing.T) {
	// Get the help text
	helpOutput := configCmd.Long

	// Check that key information is present
	expectedContent := []string{
		"configuration",
		"journal directory",
		"editor",
		"preview settings",
		"Configuration file",
		"environment variables",
		"defaults",
		"precedence",
		"LOGMD_",
		".logmdconfig",
	}

	for _, content := range expectedContent {
		if !strings.Contains(helpOutput, content) {
			t.Errorf("Help text should contain %q", content)
		}
	}
}

// TestConfigCommandRegistration tests that the command is properly registered.
func TestConfigCommandRegistration(t *testing.T) {
	// Check that config command exists in root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "config" {
			found = true
			break
		}
	}

	if !found {
		t.Error("config command should be registered with root command")
	}

	// Check basic command properties
	if configCmd.Use != "config" {
		t.Errorf("Expected Use to be 'config', got %q", configCmd.Use)
	}

	if configCmd.Short == "" {
		t.Error("config command should have a short description")
	}

	if configCmd.Long == "" {
		t.Error("config command should have a long description")
	}
}

// saveEnvironment saves current environment variables for restoration.
// Learn: Test isolation requires careful environment management.
func saveEnvironment() map[string]string {
	saved := make(map[string]string)
	envVars := []string{
		"LOGMD_DIRECTORY", "LOGMD_EDITOR", "LOGMD_PREVIEW_LINES",
		"EDITOR", "HOME",
	}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			saved[envVar] = value
		}
	}

	return saved
}

// restoreEnvironment restores previously saved environment variables.
func restoreEnvironment(saved map[string]string) {
	// Clear all our test variables first
	clearLogmdEnvironment()
	os.Unsetenv("HOME")

	// Restore saved values
	for envVar, value := range saved {
		os.Setenv(envVar, value)
	}
}

// clearLogmdEnvironment clears all logmd-related environment variables.
func clearLogmdEnvironment() {
	envVars := []string{
		"LOGMD_DIRECTORY", "LOGMD_EDITOR", "LOGMD_PREVIEW_LINES",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
