// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package gzip

import (
	"compress/gzip"
	"io"
	"os"
)

func GzipFile(source, target string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := gzip.NewWriter(out)
	defer writer.Close()

	_, err = io.Copy(writer, in)
	return err
}

func GunzipFile(source, target string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	reader, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer reader.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	return err
}
