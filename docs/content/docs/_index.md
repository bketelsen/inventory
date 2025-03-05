---
title: 'Documentation'
---

Inventory is a single-binary application that includes a server, a client, and a CLI. 



## The Server
The Inventory server runs an RPC service on port `9999` by default. It accepts inventory reports from one or more agents running on your servers.

Inventory reports are stored in memory and are not persistent. Each time an agent sends a report it replaces the previous report stored for that agent.

Start a server manually by running `inventory serve` in a terminal:

```bash
$ inventory serve
```

Inventory connects to the APIs of supported container and VM runtimes directly, so the user running the `inventory serve` command should belong
to the appropriate groups to read container/VM status information.  Typically this is the `docker` group for Docker and the `incus-admin` group for Incus.

The server should be run as a persistent, long-running service. [SystemD example units are provided in the GitHub repository](https://github.com/bketelsen/inventory/blob/main/contrib/inventory-server.service).

Inventory also runs an HTTP server on port `8000` by default. It serves a dashboard with details about your deployments.

![dashboard](/images/dashboard.png "Inventory Dashboard")

The Inventory server can be started without a [configuration](#configuration) file in place, and will listen on the default ports (9999/8000).

## The Reporter
The Inventory reporter collects information about services, containers, and network listeners running locally and sends them via RPC to the Inventory server.

Container details for Docker and Incus are collected automatically.

You can manually specify services you want to track by adding them to the [configuration file](#configuration). This is useful for services that are installed directly on the host, not in a container runtime. 

The Inventory reporter can not be run without a [configuration](#configuration) file in place.

## The CLI
You can use the `inventory search` command from a terminal to search for services, containers and listeners.

```bash
$ inventory search caddy
```

```
Found 1 server with "caddy"

Found 3 Listeners on selfie / 10.0.1.46:
LISTEN ADDRESS  PORT   PID  PROGRAM  
127.0.0.1       2019  1557  caddy    
::               443  1557  caddy    
::                80  1557  caddy    

Found 1 Services on selfie / 10.0.1.46:
NAME   PORT  LISTEN ADDRESS  PROTOCOL  UNIT           
caddy   443  10.0.0.46       tcp       caddy.service  
```

The Inventory CLI can not be run without a [configuration](#configuration) file in place.

## Configuration
Inventory's configuration file is in YAML format and must be named `inventory` - *no ".yaml" extension*.

Inventory searches for the `inventory` file in the following directories:

1. `/etc/inventory/`
2. `$HOME/.inventory/`
3. `$PWD`               

### Server Configuration Keys
```yaml
http_port: 8000 # Dashboard HTTP port
rpc_port: 9999  # RPC Service port
```

### Client Configuration Keys
```yaml
server:
  address: 192.168.1.10:9999 # address/port of the collector
location: Rack 5, Slot 6     # free-form text to help locate the server
description: 2U AMD w/NVIDIA # free-form description of the server
services: #see below

```

The `location` and `description` values are displayed on the Dashboard.

If you want to track services running on the host, you can list them in the `services` key:

```yaml
services:
    - name: syncthing               # free-form name
      port: 8384                    # listen port
      listenaddress: 192.168.1.165  # listen address
      protocol: tcp                 # listen protocol
      unit: syncthing@.service      # systemd unit
```


### Common Configuration Keys

```yaml
verbose: true
```