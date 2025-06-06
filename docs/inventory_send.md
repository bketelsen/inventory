# inventory send

Send host and container information to the server

## Synopsis

Send host and container information to the server
This command collects information about the host and docker/incus containers
and sends it to the server. It is designed to be run as a cron job or systemd timer.
It is not intended to be run interactively.
	

```
inventory send [flags]
```

## Examples

```
inventory send

// more verbose output
inventory send --verbose

// specify a config file
inventory send --verbose --config /path/to/config.yaml
```

## Options

```
  -d, --description string   Description of the server (default "My Generic Server")
  -h, --help                 help for send
  -l, --location string      Location of the server (default "My Home Lab")
  -r, --remote string        Remote inventory server address (default "10.0.1.1:9999")
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/home/bjk/inventory/inventory.yml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory](inventory.md)	 - Inventory is a tool to collect and report deployment information

