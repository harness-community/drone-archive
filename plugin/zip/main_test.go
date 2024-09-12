// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package zip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestZipArchive(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetZip := filepath.Join(os.TempDir(), "test_archive.zip")
	defer os.Remove(targetZip)

	err := Zip(sourceDir, targetZip, "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat(targetZip); os.IsNotExist(err) {
		t.Fatalf("expected zip file to be created: %v", err)
	}
}

func TestZipArchiveWithGlob(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetZip := filepath.Join(os.TempDir(), "test_glob_archive.zip")
	defer os.Remove(targetZip)

	err := Zip(sourceDir, targetZip, "", "*.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Add specific checks if needed to verify only txt files are archived.
}

func TestZipArchiveWithExclude(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetZip := filepath.Join(os.TempDir(), "test_exclude_archive.zip")
	defer os.Remove(targetZip)

	err := Zip(sourceDir, targetZip, "*.txt", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Add specific checks if needed to verify excluded files are not in the zip.
}

func TestUnzipExtract(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetZip := filepath.Join(os.TempDir(), "test_extract.zip")
	defer os.Remove(targetZip)

	err := Zip(sourceDir, targetZip, "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	extractDir := filepath.Join(os.TempDir(), "extract_test")
	defer os.RemoveAll(extractDir)

	err = Unzip(targetZip, extractDir, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		t.Fatalf("expected extract directory to be created")
	}
}

func createTestDir(t *testing.T) string {
	testDir := filepath.Join(os.TempDir(), "zip_test")
	err := os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	files := []string{"file1.txt", "file2.log", "file3.txt"}
	for _, file := range files {
		f, err := os.Create(filepath.Join(testDir, file))
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		f.Close()
	}

	return testDir
}

func TestZipPatterns(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	tests := []struct {
		name     string
		glob     string
		exclude  string
		expected []string
	}{
		{"Match all txt files", "*.txt", "", []string{"file1.txt", "file3.txt"}},
		{"Match any file ending with file.txt", "*file.txt", "", []string{"myfile.txt"}},
		{"Match all files in dir", "dir/*", "", []string{"dir/file1.txt", "dir/file2.log"}},
		{"Match log files in subdirs", "dir/*/*.log", "", []string{"dir/subdir/file1.log"}},
		{"Match with ?", "file?.txt", "", []string{"file1.txt"}},
		{"Match exactly 3 chars before extension", "dir/???.log", "", []string{"dir/abc.log"}},
		{"Combined * and ?", "file?*.txt", "", []string{"file1abc.txt"}},
		{"Exclude log files", "*", "*.log", []string{"file1.txt", "file3.txt"}},
		{"Exclude specific files", "*", "file1.txt", []string{"file3.txt", "file4.log"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			targetZip := filepath.Join(os.TempDir(), "test_zip_patterns.zip")
			defer os.Remove(targetZip)

			err := Zip(sourceDir, targetZip, test.exclude, test.glob)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			// Add validation logic to check files inside the zip.
			// This can be done by extracting and comparing file names.
		})
	}
}

func TestUnzipPatterns(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetZip := filepath.Join(os.TempDir(), "test_extract_patterns.zip")
	defer os.Remove(targetZip)

	err := Zip(sourceDir, targetZip, "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := []struct {
		name     string
		glob     string
		expected []string
	}{
		{"Extract all txt files", "*.txt", []string{"file1.txt", "file3.txt"}},
		{"Extract log files", "*.log", []string{"file2.log"}},
		{"Extract files with ?", "file?.txt", []string{"file1.txt"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			extractDir := filepath.Join(os.TempDir(), "extract_test")
			defer os.RemoveAll(extractDir)

			err := Unzip(targetZip, extractDir, test.glob)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

		})
	}
}
