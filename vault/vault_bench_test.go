package vault

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// BenchmarkListEntries tests performance with large numbers of entries.
// Learn: Benchmark functions start with "Benchmark" and take *testing.B.
// See: https://pkg.go.dev/testing#hdr-Benchmarks
func BenchmarkListEntries(b *testing.B) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "logmd-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		b.Fatalf("New() failed: %v", err)
	}

	// Create test files (simulate 5000 entries)
	entryCount := 5000
	baseDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < entryCount; i++ {
		date := baseDate.AddDate(0, 0, i)
		filename := date.Format("2006-01-02.md")
		path := vault.Directory + "/" + filename
		if err := os.WriteFile(path, []byte("# "+date.Format("2006-01-02")+"\n\nContent"), 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		entries, err := vault.ListEntries()
		if err != nil {
			b.Fatalf("ListEntries() failed: %v", err)
		}
		if len(entries) != entryCount {
			b.Fatalf("Expected %d entries, got %d", entryCount, len(entries))
		}
	}
}

// BenchmarkListEntriesInfo tests performance of metadata retrieval.
func BenchmarkListEntriesInfo(b *testing.B) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "logmd-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		b.Fatalf("New() failed: %v", err)
	}

	// Create test files (smaller set for metadata testing)
	entryCount := 1000
	baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < entryCount; i++ {
		date := baseDate.AddDate(0, 0, i)
		filename := date.Format("2006-01-02.md")
		path := vault.Directory + "/" + filename
		content := fmt.Sprintf("# %s\n\nEntry %d content", date.Format("2006-01-02"), i)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		entries, err := vault.ListEntriesInfo()
		if err != nil {
			b.Fatalf("ListEntriesInfo() failed: %v", err)
		}
		if len(entries) != entryCount {
			b.Fatalf("Expected %d entries, got %d", entryCount, len(entries))
		}
	}
}

// BenchmarkCreateEntry tests entry creation performance.
func BenchmarkCreateEntry(b *testing.B) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "logmd-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		b.Fatalf("New() failed: %v", err)
	}

	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		date := baseDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		err := vault.CreateEntry(dateStr)
		if err != nil {
			b.Fatalf("CreateEntry() failed: %v", err)
		}
	}
}

// BenchmarkReadEntry tests entry reading performance.
func BenchmarkReadEntry(b *testing.B) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "logmd-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		b.Fatalf("New() failed: %v", err)
	}

	// Create a single test entry
	testDate := "2024-01-15"
	testContent := "# 2024-01-15\n\nThis is a test entry with some content to read."
	err = vault.WriteEntry(testDate, []byte(testContent))
	if err != nil {
		b.Fatalf("Failed to create test entry: %v", err)
	}

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		content, err := vault.ReadEntry(testDate)
		if err != nil {
			b.Fatalf("ReadEntry() failed: %v", err)
		}
		if len(content) == 0 {
			b.Fatal("Expected content, got empty")
		}
	}
}

// BenchmarkEntryExists tests existence checking performance.
func BenchmarkEntryExists(b *testing.B) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "logmd-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vault, err := New(tmpDir)
	if err != nil {
		b.Fatalf("New() failed: %v", err)
	}

	// Create a test entry
	testDate := "2024-01-15"
	err = vault.CreateEntry(testDate)
	if err != nil {
		b.Fatalf("Failed to create test entry: %v", err)
	}

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		exists := vault.EntryExists(testDate)
		if !exists {
			b.Fatal("Entry should exist")
		}
	}
}
