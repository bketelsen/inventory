# Installing Inventory

You can install inventory by downloading a release from GitHub or by using our installer script.

Choose your adventure below.

## Server Only / Docker Compose

To run just the server in Docker compose, use the `compose.yaml` in the `/contrib` folder as a reference:

```yaml
services:
  inventory:
    image: ghcr.io/bketelsen/inventory:0.7.3
    ports:
      - 9999:9999
      - 8000:8000
    restart: always
    command: ["server"]
    volumes:
      - type: bind
        source: ./inventory.yaml
        target: /etc/inventory/inventory.yaml
        read_only: true
```

## Direct Download

You can download the binary from the [inventory releases page](https://github.com/bketelsen/inventory/releases) on GitHub and add to your `$PATH`.

The inventory_VERSION_checksums.txt file contains the SHA-256 checksum for each file.

## Installer Script

We also have an [install script](https://github.com/bketelsen/inventory/blob/main/install.sh) which is very useful in scenarios like CI.

By default, it installs on the `./bin` directory relative to the working directory:

```bash
sh -c "$(curl --location https://bketelsen.github.io/inventory/install.sh)" -- -d
```

It is possible to override the installation directory with the -b parameter. On Linux, common choices are `~/.local/bin` and `~/bin` to install for the current user or `/usr/local/bin` to install for all users:

```bash
sh -c "$(curl --location https://bketelsen.github.io/inventory/install.sh)" -- -d -b ~/.local/bin
```

!> On macOS and Windows, ~/.local/bin and ~/bin are not added to $PATH by default.

By default, it installs the latest version available. You can also specify a tag ([available in releases](https://github.com/bketelsen/inventory/releases)) to install a specific version:

```bash
sh -c "$(curl --location https://bketelsen.github.io/inventory/install.sh)" -- -d v0.2.2
```