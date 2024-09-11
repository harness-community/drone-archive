// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"github.com/harness-community/drone-archive/plugin/gzip"
	"github.com/harness-community/drone-archive/plugin/tar"
	"github.com/harness-community/drone-archive/plugin/zip"
	"os"
)

type Plugin struct {
	Source      string `envconfig:"PLUGIN_SOURCE"`
	Target      string `envconfig:"PLUGIN_TARGET"`
	Format      string `envconfig:"PLUGIN_FORMAT"`
	Action      string `envconfig:"PLUGIN_ACTION"` // "archive" or "extract"
	Overwrite   bool   `envconfig:"PLUGIN_OVERWRITE"`
	TarCompress bool   `envconfig:"PLUGIN_TARCOMPRESS"`
	Glob        string `envconfig:"PLUGIN_GLOB"`
	LogLevel    string `envconfig:"PLUGIN_LOG_LEVEL"`
}

func (p *Plugin) Exec(ctx context.Context) error {
	if p.Overwrite {
		if _, err := os.Stat(p.Target); err == nil {
			if err := os.RemoveAll(p.Target); err != nil {
				return fmt.Errorf("failed to overwrite target: %w", err)
			}
		}
	}

	switch p.Format {
	case "zip":
		if p.Action == "archive" {
			return zip.Zip(p.Source, p.Target, p.Glob)
		} else {
			return zip.Unzip(p.Source, p.Target, p.Glob)
		}
	case "tar":
		return p.handleTarAction()
	case "gzip":
		if p.Action == "archive" {
			return gzip.GzipFile(p.Source, p.Target)
		} else {
			return gzip.GunzipFile(p.Source, p.Target)
		}
	default:
		return fmt.Errorf("unsupported format: %s", p.Format)
	}
}

func (p *Plugin) handleTarAction() error {
	if p.Action == "archive" {
		return tar.Tar(p.Source, p.Target, p.Glob, p.TarCompress)
	}
	return tar.Untar(p.Source, p.Target, p.Glob)
}
