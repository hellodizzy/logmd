// Package vault provides file system operations for the logmd journal directory.
// This package handles directory creation, file enumeration, and path
// resolution for daily journal entries stored as markdown files.
//
// Learn: Package names should be short, clear, and lowercase without underscores.
// See: https://go.dev/blog/package-names
package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Vault represents a journal directory with its path and configuration.
// Learn: Struct types should have clear, descriptive names and documented fields.
// See: https://go.dev/doc/effective_go#names
type Vault struct {
	// Directory is the absolute path to the journal's root directory
	Directory string
}

// EntryInfo contains metadata about a journal entry.
// Learn: Structs can contain various data types to group related information.
// See: https://go.dev/tour/moretypes/2
type EntryInfo struct {
	// Date is the entry date in YYYY-MM-DD format
	Date string
	// Path is the absolute file path to the entry
	Path string
	// Exists indicates whether the file exists on disk
	Exists bool
	// Size is the file size in bytes (0 if file doesn't exist)
	Size int64
	// ModTime is the last modification time
	ModTime time.Time
}

// New creates a new Vault instance with the given directory path.
// It ensures the directory exists with proper permissions (0700).
// Learn: Constructor functions in Go typically start with "New" and return pointers.
// See: https://go.dev/doc/effective_go#constructors
func New(directory string) (*Vault, error) {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(absDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", absDir, err)
	}

	return &Vault{Directory: absDir}, nil
}

// TodayPath returns the file path for today's journal entry.
// The filename follows the format YYYY-MM-DD.md using local timezone.
// Learn: Methods in Go are functions with receiver arguments.
// See: https://go.dev/tour/methods/1
func (v *Vault) TodayPath() string {
	today := time.Now().Format("2006-01-02")
	return filepath.Join(v.Directory, today+".md")
}

// DatePath returns the file path for a specific date's journal entry.
// The date string must be in YYYY-MM-DD format.
func (v *Vault) DatePath(date string) string {
	return filepath.Join(v.Directory, date+".md")
}

// EntryExists checks if a journal entry exists for the given date.
// Learn: Boolean functions should clearly indicate what they're checking.
// See: https://go.dev/doc/effective_go#names
func (v *Vault) EntryExists(date string) bool {
	path := v.DatePath(date)
	_, err := os.Stat(path)
	return err == nil
}

// TodayExists checks if today's journal entry exists.
func (v *Vault) TodayExists() bool {
	today := time.Now().Format("2006-01-02")
	return v.EntryExists(today)
}

// ReadEntry reads the content of a journal entry for the given date.
// Returns an error if the file doesn't exist or can't be read.
// Learn: File I/O operations should always handle errors properly.
// See: https://go.dev/doc/effective_go#errors
func (v *Vault) ReadEntry(date string) ([]byte, error) {
	path := v.DatePath(date)
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("entry %s does not exist", date)
		}
		return nil, fmt.Errorf("failed to read entry %s: %w", date, err)
	}
	return content, nil
}

// WriteEntry writes content to a journal entry for the given date.
// Creates the file if it doesn't exist, overwrites if it does.
func (v *Vault) WriteEntry(date string, content []byte) error {
	path := v.DatePath(date)
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write entry %s: %w", date, err)
	}
	return nil
}

// CreateEntry creates a new journal entry with a simple template.
// Returns an error if the file already exists.
func (v *Vault) CreateEntry(date string) error {
	if v.EntryExists(date) {
		return fmt.Errorf("entry %s already exists", date)
	}

	template := fmt.Sprintf("# %s\n\n", date)
	return v.WriteEntry(date, []byte(template))
}

// CreateTodayEntry creates today's journal entry with a simple template.
// Returns an error if today's entry already exists.
func (v *Vault) CreateTodayEntry() error {
	today := time.Now().Format("2006-01-02")
	return v.CreateEntry(today)
}

// GetEntryInfo returns metadata about a journal entry.
// Learn: Methods can return structs to group related information.
func (v *Vault) GetEntryInfo(date string) EntryInfo {
	path := v.DatePath(date)
	info := EntryInfo{
		Date:   date,
		Path:   path,
		Exists: false,
		Size:   0,
	}

	if stat, err := os.Stat(path); err == nil {
		info.Exists = true
		info.Size = stat.Size()
		info.ModTime = stat.ModTime()
	}

	return info
}

// ListEntries returns all journal entries sorted by date (newest first).
// Only returns .md files that match the YYYY-MM-DD.md pattern.
// Learn: Slices in Go are dynamic arrays with length and capacity.
// See: https://go.dev/blog/slices-intro
func (v *Vault) ListEntries() ([]string, error) {
	entries, err := os.ReadDir(v.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", v.Directory, err)
	}

	var mdFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".md") && isValidDateFormat(name) {
			mdFiles = append(mdFiles, name)
		}
	}

	// Sort newest first
	sort.Slice(mdFiles, func(i, j int) bool {
		return mdFiles[i] > mdFiles[j]
	})

	return mdFiles, nil
}

// ListEntriesInfo returns metadata for all journal entries sorted by date (newest first).
// This includes both existing and non-existing entries for comprehensive listing.
func (v *Vault) ListEntriesInfo() ([]EntryInfo, error) {
	filenames, err := v.ListEntries()
	if err != nil {
		return nil, err
	}

	entries := make([]EntryInfo, 0, len(filenames))
	for _, filename := range filenames {
		date := strings.TrimSuffix(filename, ".md")
		entries = append(entries, v.GetEntryInfo(date))
	}

	return entries, nil
}

// isValidDateFormat checks if filename matches YYYY-MM-DD.md pattern.
// Learn: Helper functions should be unexported (lowercase) when used only within the package.
// See: https://go.dev/doc/effective_go#names
func isValidDateFormat(filename string) bool {
	if !strings.HasSuffix(filename, ".md") {
		return false
	}
	datePart := strings.TrimSuffix(filename, ".md")
	_, err := time.Parse("2006-01-02", datePart)
	return err == nil
}
