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
	"strings"
)

type Plugin struct {
	Source      string `envconfig:"PLUGIN_SOURCE"`
	Target      string `envconfig:"PLUGIN_TARGET"`
	Format      string `envconfig:"PLUGIN_FORMAT"`
	Action      string `envconfig:"PLUGIN_ACTION"` // "archive" or "extract"
	Overwrite   bool   `envconfig:"PLUGIN_OVERWRITE"`
	TarCompress bool   `envconfig:"PLUGIN_TARCOMPRESS"`
	Exclude     string `envconfig:"PLUGIN_EXCLUDE"`
	Glob        string `envconfig:"PLUGIN_GLOB"`
	LogLevel    string `envconfig:"PLUGIN_LOG_LEVEL"`
}

func (p *Plugin) Exec(ctx context.Context) error {
	if !p.Overwrite {
		if _, err := os.Stat(p.Target); err == nil {
			return fmt.Errorf("target file or directory already exists: %s", p.Target)
		}
	}

	switch strings.ToLower(p.Format) {
	case "zip":
		return p.handleZip()
	case "tar":
		return p.handleTar()
	case "gzip":
		return p.handleGzip()
	default:
		return fmt.Errorf("unsupported format: %s", p.Format)
	}
}

func (p *Plugin) handleZip() error {
	if strings.ToLower(p.Action) == "archive" {
		return zip.Zip(p.Source, p.Target, p.Exclude, p.Glob)
	} else if strings.ToLower(p.Action) == "extract" {
		return zip.Unzip(p.Source, p.Target, p.Glob)
	} else {
		return fmt.Errorf("unsupported action for zip: %s", p.Action)
	}
}

func (p *Plugin) handleTar() error {
	if strings.ToLower(p.Action) == "archive" {
		return tar.Tar(p.Source, p.Target, p.Exclude, p.Glob, p.TarCompress)
	} else if strings.ToLower(p.Action) == "extract" {
		return tar.Untar(p.Source, p.Target, p.Glob)
	} else {
		return fmt.Errorf("unsupported action for tar: %s", p.Action)
	}
}

func (p *Plugin) handleGzip() error {
	if strings.ToLower(p.Action) == "archive" {
		return gzip.GzipFile(p.Source, p.Target)
	} else if strings.ToLower(p.Action) == "extract" {
		return gzip.GunzipFile(p.Source, p.Target)
	} else {
		return fmt.Errorf("unsupported action for gzip: %s", p.Action)
	}
}
