package assist

import (
	"testing"
)

// TestMockEngine verifies that the mock engine interface implementation works.
// Learn: Interface testing ensures implementations satisfy contracts correctly.
// See: https://go.dev/doc/effective_go#interfaces
func TestMockEngine(t *testing.T) {
	engine := &MockEngine{}

	// Test that MockEngine implements Engine interface
	var _ Engine = engine

	suggestions, err := engine.Suggest("/path/to/test.md")
	if err != nil {
		t.Fatalf("Suggest() returned error: %v", err)
	}

	if len(suggestions) == 0 {
		t.Error("Expected suggestions, got empty slice")
	}

	// Verify we get the expected mock suggestions
	expectedCount := 3
	if len(suggestions) != expectedCount {
		t.Errorf("Expected %d suggestions, got %d", expectedCount, len(suggestions))
	}

	// Verify suggestions are not empty
	for i, suggestion := range suggestions {
		if suggestion == "" {
			t.Errorf("Suggestion %d is empty", i)
		}
	}
}

// TestAssistCmdExists verifies that the assist command is properly configured.
func TestAssistCmdExists(t *testing.T) {
	if AssistCmd == nil {
		t.Fatal("AssistCmd should not be nil")
	}

	if AssistCmd.Use != "assist" {
		t.Errorf("Expected Use='assist', got Use='%s'", AssistCmd.Use)
	}

	if AssistCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if AssistCmd.Long == "" {
		t.Error("Long description should not be empty")
	}

	if AssistCmd.Run == nil {
		t.Error("Run function should not be nil")
	}
}
