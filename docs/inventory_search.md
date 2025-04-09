# inventory search

search for services or containers

## Synopsis

Search returns a list of services, listeners, and containers that match the query.

```
inventory search [query] [flags]
```

## Examples

```
inventory search jellyfin

// more verbose output
inventory search --verbose jellyfin

// specify a config file
inventory search --verbose --config /path/to/config.yaml jellyfin
```

## Options

```
  -h, --help            help for search
  -r, --remote string   Remote inventory server address (default "10.0.1.1:9999")
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/var/home/bjk/projects/inventory/inventory.yml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory](inventory.md)	 - Inventory is a tool to collect and report deployment information

