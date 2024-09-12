// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package tar

import (
	"archive/tar"
	"compress/gzip"
	"github.com/bmatcuk/doublestar/v4"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file
	if strings.HasSuffix(source, ".gz") {
		reader, err = gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer reader.(*gzip.Reader).Close()
	}

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		matchesGlob, _ := doublestar.Match(globPattern, header.Name)

		if globPattern != "" && !matchesGlob {
			// Skip this file
			continue
		}

		targetPath := filepath.Join(target, header.Name)

		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(targetPath, 0755)
			continue
		}

		if header.Typeflag == tar.TypeReg {
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
