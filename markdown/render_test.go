package markdown

import (
	"strings"
	"testing"
)

// TestNewRenderer tests the renderer constructor.
// Learn: Constructor tests should verify the object is properly initialized.
func TestNewRenderer(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	if renderer == nil {
		t.Error("NewRenderer() returned nil renderer")
	}

	if renderer.glamourRenderer == nil {
		t.Error("glamourRenderer should not be nil")
	}

	if renderer.goldmarkParser == nil {
		t.Error("goldmarkParser should not be nil")
	}
}

// TestRenderBasicMarkdown tests rendering of basic markdown elements.
func TestRenderBasicMarkdown(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		contains []string // Strings that should be present in output
	}{
		{
			name:     "SimpleHeading",
			input:    "# Main Title\n\nSome content here.",
			contains: []string{"Main Title"},
		},
		{
			name:     "BoldText",
			input:    "This is **bold text** in a sentence.",
			contains: []string{"bold text"},
		},
		{
			name:     "ItalicText",
			input:    "This is *italic text* in a sentence.",
			contains: []string{"italic text"},
		},
		{
			name:     "CodeBlock",
			input:    "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
			contains: []string{"func main", "Hello"},
		},
		{
			name:     "InlineCode",
			input:    "Use the `fmt.Println()` function to print.",
			contains: []string{"fmt.Println()"},
		},
		{
			name:     "Blockquote",
			input:    "> This is a blockquote\n> with multiple lines.",
			contains: []string{"blockquote", "multiple lines"},
		},
		{
			name:     "UnorderedList",
			input:    "- First item\n- Second item\n- Third item",
			contains: []string{"First item", "Second item", "Third item"},
		},
		{
			name:     "OrderedList",
			input:    "1. First step\n2. Second step\n3. Third step",
			contains: []string{"First step", "Second step", "Third step"},
		},
		{
			name:     "Table",
			input:    "| Name | Age |\n|------|-----|\n| John | 25 |\n| Jane | 30 |",
			contains: []string{"Name", "Age", "John", "Jane"},
		},
		{
			name:     "Strikethrough",
			input:    "This is ~~deleted text~~ and this is normal.",
			contains: []string{"deleted text", "normal"},
		},
		{
			name:     "TaskList",
			input:    "- [x] Completed task\n- [ ] Incomplete task",
			contains: []string{"Completed task", "Incomplete task"},
		},
		{
			name:     "MultipleHeadings",
			input:    "# H1\n## H2\n### H3\n#### H4",
			contains: []string{"H1", "H2", "H3", "H4"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := renderer.Render([]byte(tc.input))
			if err != nil {
				t.Fatalf("Render() failed: %v", err)
			}

			if result == "" {
				t.Error("Render() returned empty string")
			}

			// Check that expected content is present
			for _, expected := range tc.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Rendered output should contain %q, but got:\n%s", expected, result)
				}
			}
		})
	}
}

// TestRenderComplexMarkdown tests rendering of complex markdown documents.
func TestRenderComplexMarkdown(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	complexMarkdown := `# Daily Journal Entry

## Morning Reflections

Today started with a **beautiful sunrise** and *gentle breeze*. I feel:

- Energized and ready to tackle the day
- Grateful for the peaceful morning
- Excited about the projects ahead

## Technical Work

I've been working on a CLI tool with the following features:

1. **Vault management** - File system operations
2. **Timeline interface** - Interactive TUI with Bubble Tea
3. **Markdown rendering** - Beautiful terminal display

Here's a code snippet I wrote:

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, logmd!\")\n}\n```" + `

> "The best way to predict the future is to create it." - Peter Drucker

## Task Progress

- [x] Implement vault package
- [x] Create timeline TUI  
- [ ] Add view command
- [ ] Write comprehensive tests

## Data Summary

| Feature | Status | Priority |
|---------|--------|----------|
| Vault   | ✅ Done | High     |
| Timeline| ✅ Done | High     |
| View    | 🚧 WIP  | Medium   |

## Notes

Some things I learned today:
- Go's ~~interface system~~ type system is elegant
- Bubble Tea makes TUI development fun
- Testing is crucial for reliability

That's all for today! 🚀`

	result, err := renderer.Render([]byte(complexMarkdown))
	if err != nil {
		t.Fatalf("Render() failed for complex markdown: %v", err)
	}

	if result == "" {
		t.Error("Render() returned empty string for complex markdown")
	}

	// Check for key elements
	expectedElements := []string{
		"Daily Journal Entry",
		"Morning Reflections",
		"beautiful sunrise",
		"gentle breeze",
		"Energized and ready",
		"Technical Work",
		"CLI tool",
		"Hello, logmd!",
		"Peter Drucker",
		"Task Progress",
		"Data Summary",
		"Feature",
		"Status",
		"Priority",
		"interface system",
		"Bubble Tea",
		"🚀",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("Complex markdown should contain %q", expected)
		}
	}
}

// TestRenderEmptyContent tests rendering of empty or whitespace-only content.
func TestRenderEmptyContent(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "EmptyString",
			input: "",
		},
		{
			name:  "WhitespaceOnly",
			input: "   \n\n\t\t\n   ",
		},
		{
			name:  "NewlinesOnly",
			input: "\n\n\n\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := renderer.Render([]byte(tc.input))
			if err != nil {
				t.Fatalf("Render() failed for %s: %v", tc.name, err)
			}

			// Result should be non-nil (even if empty content)
			// The exact behavior depends on glamour's handling of empty content
			_ = result // We don't assert specific behavior for empty content
		})
	}
}

// TestRenderWithSpecialCharacters tests rendering with special characters.
func TestRenderWithSpecialCharacters(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	specialContent := `# Special Characters Test

Unicode symbols: 🎉 🚀 ✅ ❌ ⭐ 💡 📝 🔍

Math symbols: α β γ δ ε ∑ ∏ ∆ ∇ ∞

Quotes: "Smart quotes" and 'single quotes'

Dashes: em-dash — and en-dash –

Arrows: ← → ↑ ↓ ↔ ↕

Special punctuation: ¡Hola! ¿Cómo estás?`

	result, err := renderer.Render([]byte(specialContent))
	if err != nil {
		t.Fatalf("Render() failed for special characters: %v", err)
	}

	// Check that the content was processed without errors
	if result == "" {
		t.Error("Render() returned empty string for special characters")
	}

	// Check for some key elements
	expectedElements := []string{
		"Special Characters Test",
		"🎉", "🚀", "✅",
		"α", "β", "∑",
		"Smart quotes",
		"em-dash",
		"¡Hola!",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("Special characters test should contain %q", expected)
		}
	}
}

// TestRenderJournalEntryFormat tests rendering of typical journal entry format.
func TestRenderJournalEntryFormat(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	journalEntry := `# 2024-01-15

## Morning

Woke up feeling refreshed after a good night's sleep. The weather looks promising for a productive day.

**Goals for today:**
- Complete the view command implementation
- Write comprehensive tests
- Update documentation

## Work Session

Implemented the markdown renderer with the following features:

` + "```markdown\n# Headers\n**Bold text**\n*Italic text*\n- Lists\n```" + `

The renderer supports:
- GitHub Flavored Markdown
- Syntax highlighting for code blocks
- Tables and task lists
- Beautiful terminal styling

## Evening Reflection

> "A day without learning is a day wasted."

Today was productive. I learned more about:
1. Go's markdown ecosystem
2. Terminal rendering with glamour
3. Testing strategies for CLI tools

## Tomorrow's Plan

- [ ] Implement Phase 5 features
- [ ] Add more configuration options
- [ ] Consider adding export functionality

Weather: ☀️ Sunny  
Mood: 😊 Satisfied  
Energy: 🔋 High`

	result, err := renderer.Render([]byte(journalEntry))
	if err != nil {
		t.Fatalf("Render() failed for journal entry: %v", err)
	}

	if result == "" {
		t.Error("Render() returned empty string for journal entry")
	}

	// Check for journal-specific elements
	expectedElements := []string{
		"2024-01-15",
		"Morning",
		"Goals for today",
		"Work Session",
		"markdown renderer",
		"GitHub Flavored Markdown",
		"Evening Reflection",
		"day without learning",
		"Tomorrow's Plan",
		"☀️ Sunny",
		"😊 Satisfied",
		"🔋 High",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("Journal entry should contain %q", expected)
		}
	}
}
