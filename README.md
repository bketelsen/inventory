# Inventory

Inventory is an application that tracks deployed services/containers. It was built with a homelab in mind. It aims to answer the question:

| "Where the heck did I deploy Jellyfin, and was it Docker or Incus?"

## Installation

Put the binary somewhere you can find it. `/usr/local/bin/` never hurts.

Example systemd unit files are in the `/contrib` folder, along with example crontab and logrotate configurations.

## Configuration

`inventory` searches `/etc/inventory/`, `$HOME/.config/inventory/` and `$HOME/.inventory/`(deprecated) for a yaml formatted file named `inventory.yml` with configuration values.

Example config:

```yaml
client:
    description: Generic Server
    location: Home Lab
    remote: 192.168.5.1:9999
debug: false
log-level: 0
server:
    http-port: 8000
    listen: 0.0.0.0
    rpc-port: 9999
```

If you want to track services that aren't deployed via docker or incus, you can add them as an array to the config file like this:

```yaml
services:
    - name: syncthing
      port: 0
      listeners:
        - port: 8384
          listen_address: 0.0.0.0
          protocol: tcp
        - port: 22000
          listen_address: 0.0.0.0
          protocol: tcp
      protocol: ""
      unit: syncthing@.service
```

## Permissions

The `inventory send` command should be run by a user who belongs to the `docker` and `incus-admin` groups. This can be root,
or any other user in those groups.


## Web Template

Web template modified from [AdminLTE](https://github.com/ColorlibHQ/AdminLTE), MIT License, Copyright (c) 2014-2023 ColorlibHQ




