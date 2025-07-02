package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"logmd/config"
	"logmd/vault"
)

// todayCmd represents the today command
// Learn: Each command in Cobra is a struct that defines its behavior and flags.
// See: https://pkg.go.dev/github.com/spf13/cobra#Command
var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Open today's journal entry for editing",
	Long: `Opens today's journal entry in your preferred editor. If the entry doesn't
exist, it will be created with a simple template. The file is saved in the
configured journal directory with the format YYYY-MM-DD.md.`,
	RunE: runTodayCommand,
}

// runTodayCommand implements the core logic for the today command.
// Learn: Separating command logic into functions makes testing easier.
// See: https://go.dev/doc/effective_go#functions
func runTodayCommand(cmd *cobra.Command, args []string) error {
	// Step 1: Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Step 2: Create vault instance (handles directory creation)
	v, err := vault.New(cfg.Directory)
	if err != nil {
		return fmt.Errorf("failed to initialize journal directory: %w", err)
	}

	// Step 3: Get today's date and check if entry exists
	today := time.Now().Format("2006-01-02")
	entryPath := v.TodayPath()

	// Step 4: Create today's entry if it doesn't exist
	if !v.TodayExists() {
		err = v.CreateTodayEntry()
		if err != nil {
			return fmt.Errorf("failed to create today's entry: %w", err)
		}
		fmt.Printf("Created new journal entry: %s\n", today)
	} else {
		fmt.Printf("Opening existing journal entry: %s\n", today)
	}

	// Step 5: Launch editor
	err = launchEditor(cfg.Editor, entryPath)
	if err != nil {
		return fmt.Errorf("failed to launch editor: %w", err)
	}

	fmt.Printf("Journal entry saved: %s\n", entryPath)
	return nil
}

// launchEditor spawns the specified editor with the given file path.
// Learn: os/exec package is used to run external programs from Go.
// See: https://pkg.go.dev/os/exec#Cmd
func launchEditor(editor, filePath string) error {
	// Create command to launch editor
	cmd := exec.Command(editor, filePath)

	// Connect stdin, stdout, stderr to allow interactive editing
	// Learn: This allows the editor to interact with the user normally.
	// See: https://pkg.go.dev/os/exec#Cmd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to complete
	err := cmd.Run()
	if err != nil {
		// Check if it's an exit status error (editor exited non-zero)
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("editor exited with status %d", exitError.ExitCode())
		}
		// Other errors (command not found, permission denied, etc.)
		return fmt.Errorf("failed to run editor '%s': %w", editor, err)
	}

	return nil
}

func init() {
	// Learn: init() functions run automatically when the package is imported.
	// This is how Cobra commands are typically registered.
	rootCmd.AddCommand(todayCmd)
}
