# drone-archive

- [Synopsis](#Synopsis)
- [Plugin Image](#Plugin-Image)
- [Parameters](#Parameters)
- [Building](#building)
- [Examples](#Examples)


## Synopsis

This repository contains a plugin for running zip/unzip compression.
Currently, it supports zip/unzip functionalities.

## Plugin Image

The plugin `harnesscommunitytest/drone-archive` is available for the following architectures:

| OS            | Tag             |
|---------------|-----------------|
| latest        | `latest`        |
| linux/amd64   | `linux-amd64`   |
| linux/arm64   | `linux-arm64`   |
| windows/amd64 | `windows-amd64` |


## Parameters

| Parameter                                                   | Comments    |
|:------------------------------------------------------------|-------------|
| source <span style="font-size: 10px"><br/>`required`</span> | source path |
| target <span style="font-size: 10px"><br/>`required`</span> | target path |


## Building

Build the plugin image:

```text
./scripts/build.sh
```

## Examples

```
docker run --rm \
    -e PLUGIN_SOURCE="$SOURCE_PATH" \
    -e PLUGIN_TARGET="$TARGET_PATH" \
    drone-archive:latest

```

```
# Plugin YAML
- step:
    type: Plugin
    name: archive-plugin-arm64
    identifier: archive-plugin-arm64
    spec:
        connectorRef: harness-docker-connector
        image: harnesscommunitytest/drone-archive:linux-arm64
        settings:
            source: path/to/source
            target: targetpath
       

- step:
    type: Plugin
    name: archive-plugin-amd64
    identifier: archive-plugin-amd64
    spec:
        connectorRef: harness-docker-connector
        image: harnesscommunitytest/drone-archive:linux-amd64
        settings:
            source: path/to/source
            target: targetpath