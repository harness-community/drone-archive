package gzip

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTestFile(t *testing.T, content string) string {
	t.Helper()
	file, err := ioutil.TempFile("", "testfile_*.txt")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("unable to write to temp file: %v", err)
	}

	return file.Name()
}

func TestGzipFile(t *testing.T) {
	// Create a test file
	sourceFile := createTestFile(t, "This is a test file content")
	defer os.Remove(sourceFile)

	// Create a target gzip file
	gzipFile := filepath.Join(os.TempDir(), "testfile.gz")
	defer os.Remove(gzipFile)

	err := GzipFile(sourceFile, gzipFile)
	if err != nil {
		t.Fatalf("GzipFile() error = %v", err)
	}

	if _, err := os.Stat(gzipFile); os.IsNotExist(err) {
		t.Errorf("expected gzip file does not exist")
	}
}

func TestGunzipFile(t *testing.T) {
	// Create a test file and gzip it
	sourceFile := createTestFile(t, "This is a test file content")
	defer os.Remove(sourceFile)

	gzipFile := filepath.Join(os.TempDir(), "testfile.gz")
	defer os.Remove(gzipFile)

	err := GzipFile(sourceFile, gzipFile)
	if err != nil {
		t.Fatalf("GzipFile() error = %v", err)
	}

	unzippedFile := filepath.Join(os.TempDir(), "testfile_unzipped.txt")
	defer os.Remove(unzippedFile)

	err = GunzipFile(gzipFile, unzippedFile)
	if err != nil {
		t.Fatalf("GunzipFile() error = %v", err)
	}

	expectedContent := "This is a test file content"
	actualContent, err := ioutil.ReadFile(unzippedFile)
	if err != nil {
		t.Fatalf("unable to read unzipped file: %v", err)
	}

	if string(actualContent) != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, string(actualContent))
	}
}

func TestGzipGunzipConsistency(t *testing.T) {
	// Create a test file
	sourceFile := createTestFile(t, "Consistency check content")
	defer os.Remove(sourceFile)

	// Create gzip and gunzip files
	gzipFile := filepath.Join(os.TempDir(), "testfile.gz")
	defer os.Remove(gzipFile)

	err := GzipFile(sourceFile, gzipFile)
	if err != nil {
		t.Fatalf("GzipFile() error = %v", err)
	}

	unzippedFile := filepath.Join(os.TempDir(), "testfile_unzipped.txt")
	defer os.Remove(unzippedFile)

	err = GunzipFile(gzipFile, unzippedFile)
	if err != nil {
		t.Fatalf("GunzipFile() error = %v", err)
	}

	expectedContent := "Consistency check content"
	actualContent, err := ioutil.ReadFile(unzippedFile)
	if err != nil {
		t.Fatalf("unable to read unzipped file: %v", err)
	}

	if string(actualContent) != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, string(actualContent))
	}
}
