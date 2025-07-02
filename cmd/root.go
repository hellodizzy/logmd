package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"logmd/assist"
)

// rootCmd represents the base command when called without any subcommands
// Learn: Cobra uses a tree structure where commands can have subcommands.
// See: https://github.com/spf13/cobra/blob/main/site/content/user_guide.md
var rootCmd = &cobra.Command{
	Use:   "logmd",
	Short: "A minimal, local-first journal CLI",
	Long: `logmd is a developer-focused journaling tool that creates daily
markdown files. It provides a simple CLI interface for creating, viewing,
and browsing your daily logs.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// Learn: cobra.Execute() handles command parsing, validation, and execution flow.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Register the assist command from the assist package
	rootCmd.AddCommand(assist.AssistCmd)
}
