# inventory server

starts the RPC and HTTP servers

```
inventory server [flags]
```

## Options

```
  -h, --help            help for server
  -w, --http-port int   HTTP Port to listen on (default 8000)
  -l, --listen string   Address to listen on (default "0.0.0.0")
  -r, --rpc-port int    RPC Port to listen on (default 9999)
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/home/runner/work/inventory/inventory/inventory.yml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory](inventory.md)	 - Inventory is a tool to collect and report deployment information

