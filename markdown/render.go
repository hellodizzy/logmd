// Package markdown provides markdown parsing and ANSI rendering for logmd.
// This package uses goldmark for parsing and glamour for terminal rendering,
// ensuring beautiful display of journal entries in the terminal.
//
// Learn: Markdown processing often involves a two-step parse-then-render process.
// See: https://github.com/yuin/goldmark#overview
package markdown

import (
	"bytes"

	"github.com/charmbracelet/glamour"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Renderer handles markdown to ANSI conversion for terminal display.
// Learn: Struct types can embed behavior and state for reusable components.
// See: https://go.dev/doc/effective_go#embedding
type Renderer struct {
	glamourRenderer *glamour.TermRenderer
	goldmarkParser  goldmark.Markdown
}

// NewRenderer creates a new markdown renderer with configured styling.
// Uses glamour's auto style detection for optimal terminal appearance.
// Learn: Constructor functions should validate inputs and return configured objects.
// See: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewRenderer() (*Renderer, error) {
	// Configure glamour for terminal rendering
	glamourRenderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return nil, err
	}

	// Configure goldmark for markdown parsing
	goldmarkParser := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	return &Renderer{
		glamourRenderer: glamourRenderer,
		goldmarkParser:  goldmarkParser,
	}, nil
}

// Render converts markdown bytes to ANSI-formatted string for terminal display.
// The input should be raw markdown content read from a journal file.
// Learn: Methods that can fail should return (result, error) tuple.
// See: https://go.dev/blog/error-handling-and-go
func (r *Renderer) Render(markdown []byte) (string, error) {
	// Use glamour to render markdown with ANSI escape codes
	rendered, err := r.glamourRenderer.Render(string(markdown))
	if err != nil {
		return "", err
	}
	return rendered, nil
}

// ExtractFirstHeading parses markdown and returns the first heading after front matter.
// Returns "(untitled)" if no heading is found after YAML front matter.
// Learn: Parsing often requires state machines or careful string processing.
func ExtractFirstHeading(markdown []byte) string {
	var buf bytes.Buffer
	if err := goldmark.New().Convert(markdown, &buf); err != nil {
		return "(untitled)"
	}

	// TODO: Implement proper heading extraction after front matter
	// For Phase 0, return placeholder
	return "(untitled)"
}

// StripFrontMatter removes YAML front matter from markdown content.
// Returns the content without the leading --- delimited section.
func StripFrontMatter(content []byte) []byte {
	lines := bytes.Split(content, []byte("\n"))
	if len(lines) < 3 {
		return content
	}

	// Check for front matter delimiter
	if !bytes.Equal(lines[0], []byte("---")) {
		return content
	}

	// Find closing delimiter
	for i := 1; i < len(lines); i++ {
		if bytes.Equal(lines[i], []byte("---")) {
			// Return content after front matter
			if i+1 < len(lines) {
				return bytes.Join(lines[i+1:], []byte("\n"))
			}
			return []byte{}
		}
	}

	// No closing delimiter found, return original
	return content
}
