// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func Tar(source, target, excludePattern, globPattern string, compress bool) error {
	var fileWriter io.WriteCloser
	fileWriter, err := os.Create(target)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	var writer io.Writer = fileWriter
	if compress {
		writer = gzip.NewWriter(fileWriter)
		defer writer.(*gzip.Writer).Close()
	}

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Apply glob and exclude patterns
		matchesGlob, _ := doublestar.Match(globPattern, path)
		matchesExclude, _ := doublestar.Match(excludePattern, path)

		if (globPattern != "" && !matchesGlob) || (excludePattern != "" && matchesExclude) {
			// Skip this file or directory
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(path, source, "", -1), string(filepath.Separator))

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		return err
	})
}

func Untar(source, target, globPattern string) error {
	// Ensure the base target directory exists
	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Open the source tar file
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	var reader io.Reader = file
	// Handle .gz compressed files
	if strings.HasSuffix(source, ".gz") {
		reader, err = gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.(*gzip.Reader).Close()
	}

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %w", err)
		}

		// Match the header name with the provided glob pattern
		matchesGlob, _ := doublestar.Match(globPattern, header.Name)

		if globPattern != "" && !matchesGlob {
			// Skip this file if it doesn't match the glob pattern
			continue
		}

		// Construct the full target path for the file or directory
		targetPath := filepath.Join(target, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create the directory if it doesn't exist
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}

		case tar.TypeReg:
			// Ensure the parent directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
			}

			// Create the file and copy its content
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to copy file content to %s: %w", targetPath, err)
			}

		default:
			// Handle other file types if necessary, or skip them
			fmt.Printf("Skipping unsupported file type: %s\n", header.Name)
		}
	}

	return nil
}
