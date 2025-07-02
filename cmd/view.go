package cmd

import (
	"fmt"
	"regexp"
	"time"

	"github.com/spf13/cobra"
	"logmd/config"
	"logmd/markdown"
	"logmd/vault"
)

// viewCmd represents the view command
// Learn: Commands can accept positional arguments via the Args field or RunE function parameters.
// See: https://pkg.go.dev/github.com/spf13/cobra#PositionalArgs
var viewCmd = &cobra.Command{
	Use:   "view <YYYY-MM-DD>",
	Short: "Display a journal entry with formatted markdown",
	Long: `Renders and displays a specific journal entry using glamour for
beautiful markdown formatting. The date must match exactly the format
used for journal files (YYYY-MM-DD).

Examples:
  logmd view 2024-01-15
  logmd view 2025-06-30

The entry will be displayed with:
- Colored headings and text formatting
- Syntax-highlighted code blocks  
- Properly rendered tables and lists
- Beautiful terminal styling`,
	Args: cobra.ExactArgs(1),
	RunE: runViewCommand,
}

// runViewCommand implements the core logic for the view command.
// Learn: Separating command logic into functions makes testing and maintenance easier.
func runViewCommand(cmd *cobra.Command, args []string) error {
	dateStr := args[0]

	// Step 1: Validate date format
	if !isValidDateFormat(dateStr) {
		return fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD)", dateStr)
	}

	// Step 2: Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Step 3: Create vault instance
	v, err := vault.New(cfg.Directory)
	if err != nil {
		return fmt.Errorf("failed to initialize journal directory: %w", err)
	}

	// Step 4: Check if entry exists
	if !v.EntryExists(dateStr) {
		return fmt.Errorf("journal entry for %s does not exist", dateStr)
	}

	// Step 5: Read entry content
	content, err := v.ReadEntry(dateStr)
	if err != nil {
		return fmt.Errorf("failed to read entry %s: %w", dateStr, err)
	}

	// Step 6: Create markdown renderer
	renderer, err := markdown.NewRenderer()
	if err != nil {
		return fmt.Errorf("failed to create markdown renderer: %w", err)
	}

	// Step 7: Render and display the content
	rendered, err := renderer.Render(content)
	if err != nil {
		return fmt.Errorf("failed to render markdown: %w", err)
	}

	// Step 8: Display the rendered content
	fmt.Print(rendered)

	return nil
}

// isValidDateFormat validates that the date string matches YYYY-MM-DD format.
// Learn: Regular expressions are useful for format validation.
// See: https://pkg.go.dev/regexp
func isValidDateFormat(dateStr string) bool {
	// Check format with regex
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(dateStr) {
		return false
	}

	// Validate it's a real date
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
