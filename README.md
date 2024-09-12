# archive-plugin

- [Synopsis](#Synopsis)
- [Plugin Image](#Plugin-Image)
- [Parameters](#Parameters)
- [Building](#building)
- [Examples](#Examples)


## Synopsis

This repository contains a plugin for running archive functionalities like zip/tar/gzip.


## Plugin Image

The plugin `harnesscommunitytest/archive-plugin` is available for the following architectures:

| OS            | Tag             |
|---------------|-----------------|
| latest        | `latest`        |
| linux/amd64   | `linux-amd64`   |
| linux/arm64   | `linux-arm64`   |
| windows/amd64 | `windows-amd64` |


## Parameters

| Parameter                                                            | Comments                                                                                                                                                                  |
|:---------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| source <span style="font-size: 10px"><br/>`required`</span>          | source path                                                                                                                                                               |
| target <span style="font-size: 10px"><br/>`required`</span>          | target path                                                                                                                                                               |
| format <span style="font-size: 10px"><br/>`required`</span>          | zip/tar/gzip                                                                                                                                                              |
| action <span style="font-size: 10px"><br/>`required`</span>          | archive or extract                                                                                                                                                        |
| tarcompress <span style="font-size: 10px"><br/>`optional`</span>     | true or false (compression for tar)                                                                                                                                       |
| glob <span style="font-size: 10px"><br/>`optional`</span>            | [Ant style pattern](https://ant.apache.org/manual/dirtasks.html#patterns) of files to extract/archive from the zip/tar. Leave empty to include all files and directories. |
| exclude <span style="font-size: 10px"><br/>`optional`</span>         | [Ant style pattern](https://ant.apache.org/manual/dirtasks.html#patterns) of files to exclude from the zip/tar.                                                           |
| overwrite <span style="font-size: 10px"><br/>`optional`</span>       | true of false                                                                                                                                                             |

## Building

Build the plugin image:

```text
./scripts/build.sh
```

## Examples

```
docker run \
  -e PLUGIN_SOURCE=/data/source \
  -e PLUGIN_TARGET=/data/backup/archive.zip \
  -e PLUGIN_FORMAT=zip \
  -e PLUGIN_ACTION=archive \
  -e PLUGIN_OVERWRITE=true \
  -e PLUGIN_EXCLUDE="*.log" \
  -e PLUGIN_GLOB="**/*.txt" \
  harnesscommunitytest/archive-plugin
  
docker run \
  -e PLUGIN_SOURCE=/data/backup/archive.zip \
  -e PLUGIN_TARGET=/data/source \
  -e PLUGIN_FORMAT=zip \
  -e PLUGIN_ACTION=extract \
  -e PLUGIN_OVERWRITE=true \
  -e PLUGIN_GLOB="**/*.txt" \
  harnesscommunitytest/archive-plugin
  
docker run \
  -e PLUGIN_SOURCE=/data/source \
  -e PLUGIN_TARGET=/data/backup/archive.tar \
  -e PLUGIN_FORMAT=tar \
  -e PLUGIN_ACTION=archive \
  -e PLUGIN_OVERWRITE=true \
  -e PLUGIN_EXCLUDE="*.log" \
  -e PLUGIN_GLOB="**/*.txt" \
  -e PLUGIN_TARCOMPRESS=false \
  harnesscommunitytest/archive-plugin
  
docker run \
  -e PLUGIN_SOURCE=/data/backup/archive.tar \
  -e PLUGIN_TARGET=/data/source \
  -e PLUGIN_FORMAT=tar \
  -e PLUGIN_ACTION=extract \
  -e PLUGIN_OVERWRITE=true \
  -e PLUGIN_GLOB="**/*.txt" \
  harnesscommunitytest/archive-plugin
  
docker run \
  -e PLUGIN_SOURCE=/data/source \
  -e PLUGIN_TARGET=/data/backup/archive.tar.gz \
  -e PLUGIN_FORMAT=tar \
  -e PLUGIN_ACTION=archive \
  -e PLUGIN_OVERWRITE=true \
  -e PLUGIN_TARCOMPRESS=true \
  -e PLUGIN_EXCLUDE="*.log" \
  -e PLUGIN_GLOB="**/*.txt" \
  harnesscommunitytest/archive-plugin
  
```

```
# Plugin YAML
- step:
    type: Plugin
    name: archive-plugin-arm64
    identifier: archive-plugin-arm64
    spec:
        connectorRef: harness-docker-connector
        image: harnesscommunitytest/archive-plugin:linux-arm64
        settings:
            source: path/to/source
            target: targetpath
            format: zip/tar/gzip
            action: archive or extract
            glob: Some ant style pattern to include those files
            exclude: Some ant style pattern to exclude those files
            overwrite: true/false