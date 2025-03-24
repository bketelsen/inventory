# Inventory

Inventory is an application that tracks deployed services/containers. It was built with a homelab in mind. It aims to answer the question:

| "Where the heck did I deploy Jellyfin, and was it Docker or Incus?"

## Installation

Put the binary somewhere you can find it. `/usr/local/bin/` never hurts.

Example systemd unit files are in the `/contrib` folder, along with example crontab and logrotate configurations.

## Configuration

`inventory` searches `/etc/inventory/`, `$HOME/.config/inventory/` and `$HOME/.inventory/`(deprecated) for a yaml formatted file named `inventory.yaml` with configuration values.

Currently, you can set log level to verbose, and specify the server where reports are sent.

Example config:

```yaml
server:
  address: "10.0.1.5:9999"
verbose: false


```

If you want to track services that aren't deployed via docker or incus, you can add them as an array to the config file like this:

```yaml
services:
    - name: syncthing
      port: 8384
      listenaddress: 10.0.1.15
      protocol: tcp
      unit: syncthing@.service 

```

## Permissions

The `inventory send` command should be run by a user who belongs to the `docker` and `incus-admin` groups. This can be root,
or any other user in those groups.


## Web Template

Web template modified from [AdminLTE](https://github.com/ColorlibHQ/AdminLTE), MIT License, Copyright (c) 2014-2023 ColorlibHQ


## Available Tasks in Taskfile

### build

Description: Build the application
Summary: Build the application with ldflags to set the version with a -dev suffix.

Output: 'inventory' in project root.

Run this task:
```
task build
```

### direnv

Description: Add direnv hook to your bashrc

Run this task:
```
task direnv
```

### generate

Description: Generate CLI documentation

Run this task:
```
task generate
```

### tools

Description: Install required tools

Run this task:
```
task tools
```

### checks:all

Description: Run all go checks

Run this task:
```
task checks:all
```

### checks:format

Description: Format all Go source

Run this task:
```
task checks:format
```

### checks:staticcheck

Description: Run go staticcheck

Run this task:
```
task checks:staticcheck
```

### checks:test

Description: Run all tests

Run this task:
```
task checks:test
```

### checks:tidy

Description: Run go mod tidy

Run this task:
```
task checks:tidy
```

### checks:vet

Description: Run go vet on sources

Run this task:
```
task checks:vet
```

### docs:installer

Description: Copy installer from root to site/static directory

Run this task:
```
task docs:installer
```

### docs:site

Description: Run hugo dev server

Run this task:
```
task docs:site
```

### release:goreleaser

Description: Install goreleaser on debian derivatives

Run this task:
```
task release:goreleaser
```

### release:publish

Description: Push and tag at 0.0.1

Run this task:
```
task release:publish
```

### release:release-check

Description: Run goreleaser check

Run this task:
```
task release:release-check
```

### release:snapshot

Description: Run goreleaser in snapshot mode

Run this task:
```
task release:snapshot
```

