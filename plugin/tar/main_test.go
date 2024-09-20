// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package tar

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTarArchive(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetTar := filepath.Join(os.TempDir(), "test_archive.tar")
	defer os.Remove(targetTar)

	err := Tar(sourceDir, targetTar, "", "", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat(targetTar); os.IsNotExist(err) {
		t.Fatalf("expected tar file to be created: %v", err)
	}
}

func TestTarArchiveWithGlob(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetTar := filepath.Join(os.TempDir(), "test_glob_archive.tar")
	defer os.Remove(targetTar)

	err := Tar(sourceDir, targetTar, "", "*.txt", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

}

func TestTarArchiveWithExclude(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetTar := filepath.Join(os.TempDir(), "test_exclude_archive.tar")
	defer os.Remove(targetTar)

	err := Tar(sourceDir, targetTar, "*.txt", "", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

}

func TestUntarExtract(t *testing.T) {
	sourceDir := createTestDir(t)
	defer os.RemoveAll(sourceDir)

	targetTar := filepath.Join(os.TempDir(), "test_extract.tar")
	defer os.Remove(targetTar)

	err := Tar(sourceDir, targetTar, "", "", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	extractDir := filepath.Join(os.TempDir(), "extract_test")
	defer os.RemoveAll(extractDir)

	err = Untar(targetTar, extractDir, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		t.Fatalf("expected extract directory to be created")
	}
}

func createTestDir(t *testing.T) string {
	testDir := filepath.Join(os.TempDir(), "tar_test")
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

func TestTarPatterns(t *testing.T) {
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
			targetTar := filepath.Join(os.TempDir(), "test_tar_patterns.tar")
			defer os.Remove(targetTar)

			err := Tar(sourceDir, targetTar, test.exclude, test.glob, false)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
