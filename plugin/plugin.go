// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"github.com/harness-community/drone-archive/plugin/zip"
	"os"
	"strings"
)

type ZipPlugin struct {
	Source   string `envconfig:"PLUGIN_SOURCE"`
	Target   string `envconfig:"PLUGIN_TARGET"`
	LogLevel string `envconfig:"PLUGIN_LOG_LEVEL"`
}

func (p *ZipPlugin) Exec(ctx context.Context) error {
	sourceInfo, err := os.Stat(p.Source)
	if err != nil {
		return fmt.Errorf("error accessing source: %w", err)
	}

	if sourceInfo.IsDir() || !strings.HasSuffix(strings.ToLower(p.Source), ".zip") {
		// Zipping
		return zip.Zip(p.Source, p.Target)
	} else {
		// Unzipping
		return zip.Unzip(p.Source, p.Target)
	}
}
