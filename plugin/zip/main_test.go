// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package zip

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestZip_Exec_Zip(t *testing.T) {
	sourceDir := t.TempDir()
	targetFile := filepath.Join(t.TempDir(), "output.zip")

	file, err := os.Create(filepath.Join(sourceDir, "test.txt"))
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	file.WriteString("hello, world")
	file.Close()

	err = Zip(sourceDir, targetFile)
	if err != nil {
		t.Fatalf("failed to zip: %v", err)
	}

	_, err = zip.OpenReader(targetFile)
	if err != nil {
		t.Fatalf("output is not a valid zip file: %v", err)
	}
}

func TestZip_Exec_Unzip(t *testing.T) {
	sourceZip := filepath.Join(t.TempDir(), "test.zip")
	targetDir := t.TempDir()

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	fileWriter, err := writer.Create("test.txt")
	if err != nil {
		t.Fatalf("failed to create file in zip: %v", err)
	}
	fileWriter.Write([]byte("hello, world"))
	writer.Close()

	err = os.WriteFile(sourceZip, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("failed to write zip file: %v", err)
	}

	err = Unzip(sourceZip, targetDir)
	if err != nil {
		t.Fatalf("failed to unzip: %v", err)
	}

	extractedFile := filepath.Join(targetDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(content) != "hello, world" {
		t.Errorf("extracted content mismatch: got %v, want %v", string(content), "hello, world")
	}
}

func TestZip(t *testing.T) {
	sourceDir := t.TempDir()
	targetFile := filepath.Join(t.TempDir(), "output.zip")

	file, err := os.Create(filepath.Join(sourceDir, "test.txt"))
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	file.WriteString("dummy content")
	file.Close()

	err = Zip(sourceDir, targetFile)
	if err != nil {
		t.Fatalf("failed to zip: %v", err)
	}

	_, err = zip.OpenReader(targetFile)
	if err != nil {
		t.Fatalf("output is not a valid zip file: %v", err)
	}
}

func TestUnzip(t *testing.T) {
	sourceZip := filepath.Join(t.TempDir(), "test.zip")
	targetDir := t.TempDir()

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	fileWriter, err := writer.Create("test.txt")
	if err != nil {
		t.Fatalf("failed to create file in zip: %v", err)
	}
	fileWriter.Write([]byte("test content"))
	writer.Close()

	err = os.WriteFile(sourceZip, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("failed to write zip file: %v", err)
	}

	err = Unzip(sourceZip, targetDir)
	if err != nil {
		t.Fatalf("failed to unzip: %v", err)
	}

	extractedFile := filepath.Join(targetDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("extracted content mismatch: got %v, want %v", string(content), "test content")
	}
}