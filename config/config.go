// Package config provides configuration management for logmd using Viper.
// It handles loading settings from config files, environment variables,
// and command-line flags with a clear precedence order.
//
// Learn: Configuration packages often use the singleton pattern in Go.
// See: https://refactoring.guru/design-patterns/singleton/go/example
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration values for logmd.
// Learn: Config structs should use tag annotations for serialization.
// See: https://pkg.go.dev/reflect#StructTag
type Config struct {
	// Directory is the path to the directory where journal entries are stored
	Directory string `mapstructure:"directory"`
	// Editor is the command used to open journal files for editing
	Editor string `mapstructure:"editor"`
	// PreviewLines controls how many lines to show in timeline previews
	PreviewLines int `mapstructure:"preview_lines"`
}

// Load reads configuration from file, environment, and defaults.
// Returns a Config struct with all values resolved according to precedence.
// Learn: Viper automatically handles multiple configuration sources.
// See: https://github.com/spf13/viper#reading-config-files
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	v.SetDefault("directory", filepath.Join(homeDir, "logmd"))
	v.SetDefault("editor", getDefaultEditor())
	v.SetDefault("preview_lines", 5)

	// Configure file reading
	v.SetConfigName(".logmdconfig")
	v.SetConfigType("toml")
	v.AddConfigPath(homeDir)

	// Configure environment variables
	v.SetEnvPrefix("LOGMD")
	v.AutomaticEnv()

	// Read config file (ignore if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// getDefaultEditor returns the default editor based on environment.
// Respects $EDITOR environment variable, falls back to vim.
// Learn: Environment variable access is done through the os package.
// See: https://pkg.go.dev/os#Getenv
func getDefaultEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "vim"
}

// GetConfigPath returns the path to the configuration file.
// Returns empty string if no config file is found.
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	configPath := filepath.Join(homeDir, ".logmdconfig")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return ""
	}

	return configPath
}
