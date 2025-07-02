package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"logmd/config"
	"logmd/tui"
)

// timelineCmd represents the timeline command
// Learn: Cobra commands can be standalone or have subcommands using AddCommand().
// See: https://github.com/spf13/cobra/blob/main/site/content/user_guide.md#organizing-subcommands
var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Browse journal entries in an interactive timeline",
	Long: `Launches an interactive timeline interface using Bubble Tea TUI.
Navigate through your journal entries, expand/collapse previews, and
browse your writing history in a beautiful terminal interface.

Controls:
  ↑/k     Move up
  ↓/j     Move down
  enter   Toggle expand/collapse entry
  space   Toggle expand/collapse entry
  pgup    Page up
  pgdown  Page down
  q       Quit`,
	RunE: runTimelineCommand,
}

// runTimelineCommand implements the core logic for the timeline command.
// Learn: Separating command logic into functions makes testing and maintenance easier.
func runTimelineCommand(cmd *cobra.Command, args []string) error {
	// Step 1: Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Step 2: Create and initialize the TUI model
	model := tui.NewModel(cfg.Directory, cfg.PreviewLines)

	// Step 3: Start the Bubble Tea program
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Step 4: Run the program and handle any errors
	finalModel, err := program.Run()
	if err != nil {
		return fmt.Errorf("failed to start timeline interface: %w", err)
	}

	// Step 5: Check if the program exited with an error
	if m, ok := finalModel.(tui.Model); ok && m.Error() != nil {
		return fmt.Errorf("timeline error: %w", m.Error())
	}

	return nil
}

func init() {
	rootCmd.AddCommand(timelineCmd)
}
