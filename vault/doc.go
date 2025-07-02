/*
Package vault provides file system operations for logmd journal directories.

This package handles the core file system operations needed for managing
daily journal entries stored as markdown files. It provides functionality
for creating, reading, writing, and enumerating journal entries with proper
error handling and path resolution.

Key Features:

• Directory Management: Creates and manages journal directories with proper permissions
• Path Resolution: Generates file paths for daily entries using YYYY-MM-DD.md format
• File Operations: Read, write, and check existence of journal entries
• Entry Enumeration: List and sort journal entries by date
• Template Creation: Generate new entries with simple markdown templates
• Metadata Access: Retrieve file information including size and modification time

Usage Example:

	// Create a new vault
	vault, err := vault.New("~/journal")
	if err != nil {
		log.Fatal(err)
	}

	// Create today's entry if it doesn't exist
	if !vault.TodayExists() {
		err := vault.CreateTodayEntry()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Read an existing entry
	content, err := vault.ReadEntry("2024-01-15")
	if err != nil {
		log.Fatal(err)
	}

	// List all entries
	entries, err := vault.ListEntries()
	if err != nil {
		log.Fatal(err)
	}

File Format:

Journal entries are stored as simple markdown files with the naming convention
YYYY-MM-DD.md. New entries are created with a basic template containing just
the date as a top-level heading:

	# 2024-01-15

Directory Structure:

The vault maintains a flat directory structure where each journal entry is
a separate markdown file:

	journal/
	├── 2024-01-15.md
	├── 2024-01-14.md
	├── 2024-01-13.md
	└── ...

Error Handling:

All file operations return descriptive errors using fmt.Errorf with error
wrapping. Common error conditions include:

• Directory creation failures
• File not found errors
• Permission denied errors
• Invalid date format errors

The package uses os.MkdirTemp in tests to ensure clean, isolated test
environments without affecting the user's actual journal directory.

Learning Notes for Go Newcomers:

• File I/O: This package demonstrates proper file reading/writing with os.ReadFile and os.WriteFile
• Error Wrapping: Uses fmt.Errorf with %w verb to wrap errors while preserving the original error
• Path Manipulation: Uses filepath.Join for cross-platform path handling
• Time Formatting: Uses Go's reference time "2006-01-02" for date formatting
• Directory Permissions: Uses 0700 for directories and 0644 for files following Unix conventions
• Testing: Shows how to use os.MkdirTemp for isolated file system testing

References:

• File I/O: https://pkg.go.dev/os#ReadFile
• Error Handling: https://go.dev/blog/error-handling-and-go
• Path Manipulation: https://pkg.go.dev/path/filepath
• Time Formatting: https://go.dev/src/time/format.go
• Testing: https://pkg.go.dev/testing
*/
package vault
