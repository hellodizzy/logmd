package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"logmd/config"
)

// configCmd represents the config command
// Learn: Commands without arguments often use Run instead of RunE when no error handling is needed.
// See: https://pkg.go.dev/github.com/spf13/cobra#Command.Run
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display current configuration settings",
	Long: `Shows the active configuration including journal directory, editor,
and preview settings. Also indicates whether settings come from config
file, environment variables, or defaults.

This command helps you understand your current logmd configuration and
troubleshoot any configuration issues.

Configuration precedence (highest to lowest):
1. Environment variables (LOGMD_*)
2. Configuration file (~/.logmdconfig)  
3. Default values`,
	RunE: runConfigCommand,
}

// runConfigCommand implements the core logic for the config command.
// Learn: Separating command logic into functions makes testing and maintenance easier.
func runConfigCommand(cmd *cobra.Command, args []string) error {
	// Load current configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Display configuration information
	fmt.Println("📋 logmd Configuration")
	fmt.Println("=" + repeatString("=", 50))
	fmt.Println()

	// Show configuration file status
	configPath := config.GetConfigPath()
	if configPath != "" {
		fmt.Printf("📄 Config File: %s\n", configPath)
	} else {
		homeDir, _ := os.UserHomeDir()
		expectedPath := filepath.Join(homeDir, ".logmdconfig")
		fmt.Printf("📄 Config File: %s (not found)\n", expectedPath)
	}
	fmt.Println()

	// Display each setting with its source
	fmt.Println("⚙️  Current Settings:")
	fmt.Println()

	displaySetting("Directory", cfg.Directory, getSettingSource("LOGMD_DIRECTORY", configPath != ""))
	displaySetting("Editor", cfg.Editor, getSettingSource("LOGMD_EDITOR", configPath != ""))
	displaySetting("Preview Lines", fmt.Sprintf("%d", cfg.PreviewLines), getSettingSource("LOGMD_PREVIEW_LINES", configPath != ""))

	fmt.Println()

	// Show environment variables if set
	showEnvironmentVariables()

	// Show usage instructions
	fmt.Println("💡 Tips:")
	fmt.Printf("   • Create config file: echo 'directory = \"%s\"' > ~/.logmdconfig\n", cfg.Directory)
	fmt.Println("   • Set environment variable: export LOGMD_DIRECTORY=/path/to/journal")
	fmt.Println("   • Override editor: export LOGMD_EDITOR=code")

	return nil
}

// displaySetting shows a configuration setting with its value and source.
// Learn: Helper functions improve code readability and maintainability.
func displaySetting(name, value, source string) {
	fmt.Printf("   %-15s %s\n", name+":", value)
	fmt.Printf("   %-15s %s\n", "", source)
	fmt.Println()
}

// getSettingSource determines where a configuration setting comes from.
// Learn: Configuration source tracking helps users understand precedence.
func getSettingSource(envVar string, hasConfigFile bool) string {
	// Check if environment variable is set
	if envValue := os.Getenv(envVar); envValue != "" {
		return fmt.Sprintf("🌍 Environment variable (%s)", envVar)
	}

	// Check if we have a config file
	if hasConfigFile {
		return "📄 Configuration file (~/.logmdconfig)"
	}

	// Must be default value
	return "🔧 Default value"
}

// showEnvironmentVariables displays any set logmd environment variables.
func showEnvironmentVariables() {
	envVars := []string{"LOGMD_DIRECTORY", "LOGMD_EDITOR", "LOGMD_PREVIEW_LINES", "EDITOR"}
	hasEnvVars := false

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			if !hasEnvVars {
				fmt.Println("🌍 Environment Variables:")
				hasEnvVars = true
			}
			fmt.Printf("   %-20s %s\n", envVar+":", value)
		}
	}

	if hasEnvVars {
		fmt.Println()
	}
}

// repeatString repeats a string n times.
// Learn: Helper functions for string manipulation are common in CLI tools.
func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

func init() {
	rootCmd.AddCommand(configCmd)
}
