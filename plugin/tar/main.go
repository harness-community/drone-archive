// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package tar

import (
	"archive/tar"
	"compress/gzip"
	"github.com/bmatcuk/doublestar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Tar(source, target, globPattern string, compress bool) error {
	var fileWriter io.Writer
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	if compress {
		gzipWriter := gzip.NewWriter(tarfile)
		defer gzipWriter.Close()
		fileWriter = gzipWriter
	} else {
		fileWriter = tarfile
	}

	archive := tar.NewWriter(fileWriter)
	defer archive.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if globPattern != "" {
			match, err := doublestar.Match(globPattern, path)
			if err != nil || !match {
				return nil // Skip files not matching the glob pattern
			}
		}

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Name = strings.TrimPrefix(path, filepath.Dir(source)+"/")
		}

		if err := archive.WriteHeader(header); err != nil {
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
		_, err = io.Copy(archive, file)
		return err
	})
}

func Untar(source, target, globPattern string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	var fileReader io.Reader = file
	if strings.HasSuffix(source, ".gz") {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		fileReader = gzipReader
	}

	tarReader := tar.NewReader(fileReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)

		if globPattern != "" {
			match, err := doublestar.Match(globPattern, path)
			if err != nil || !match {
				continue // Skip files not matching the glob pattern
			}
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
		}
	}
	return nil
}
