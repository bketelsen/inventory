# Inventory

Inventory is an application that tracks deployed services/containers. It was built with a homelab in mind. It aims to answer the question:

| "Where the heck did I deploy Jellyfin, and was it Docker or Incus?"

## Installation

Put the binary somewhere you can find it. `/usr/local/bin/` never hurts.

Example systemd unit files are in the `/contrib` folder, along with example crontab and logrotate configurations.

## Configuration

`inventory` searches `/etc/inventory/` and `$HOME/.inventory` for a yaml formatted file named `inventory` with configuration values.

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
