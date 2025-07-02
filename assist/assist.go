// Package assist provides LLM-powered features for logmd.
// This package is designed to integrate AI assistance for journal writing,
// including content suggestions, writing prompts, and entry analysis.
//
// Learn: Package documentation comments should start with "Package <name>" and explain the purpose.
// See: https://go.dev/doc/effective_go#commentary
package assist

import (
	"fmt"

	"github.com/spf13/cobra"
)

// assistCmd represents the assist command (placeholder for Phase 3)
// Learn: Even placeholder code should follow Go conventions and be well-documented.
// See: https://go.dev/blog/godoc
var AssistCmd = &cobra.Command{
	Use:   "assist",
	Short: "AI-powered writing assistance (coming soon)",
	Long: `The assist command will provide AI-powered features for journal writing
including content suggestions, writing prompts, and entry analysis.
This feature is planned for Phase 3 implementation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("assist is not implemented yet. Planned for Phase 3.")
	},
}

// Engine defines the interface for LLM-powered assistance features.
// Learn: Interfaces in Go define method signatures and enable polymorphism.
// See: https://go.dev/tour/methods/9
type Engine interface {
	// Suggest generates writing suggestions based on the given file path.
	// Returns a slice of suggestion strings or an error if generation fails.
	Suggest(path string) ([]string, error)
}

// MockEngine provides a fake implementation for testing and development.
// Learn: Mock implementations are essential for testing interfaces in Go.
// See: https://go.dev/blog/testable-examples
type MockEngine struct{}

// Suggest returns hard-coded suggestions for testing purposes.
// This implementation satisfies the Engine interface for development use.
func (m *MockEngine) Suggest(path string) ([]string, error) {
	return []string{
		"Consider adding more details about your learning progress",
		"What challenges did you face today?",
		"How did you solve problems you encountered?",
	}, nil
}
