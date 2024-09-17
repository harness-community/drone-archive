// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GzipFile(source, target string) error {
	// Ensure the target directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	in, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer in.Close()

	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer out.Close()

	writer := gzip.NewWriter(out)
	defer writer.Close()

	_, err = io.Copy(writer, in)
	if err != nil {
		return fmt.Errorf("failed to compress file: %w", err)
	}

	return nil
}

func GunzipFile(source, target string) error {
	// Ensure the target directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	in, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer in.Close()

	reader, err := gzip.NewReader(in)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer reader.Close()

	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("failed to decompress file: %w", err)
	}

	return nil
}
