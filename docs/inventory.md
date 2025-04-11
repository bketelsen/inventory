# inventory

Inventory is a tool to collect and report deployment information

## Synopsis

  Inventory is a tool to collect and report deployment information to a central 
  server. It collects information about the host, docker/incus containers, and  
  manually specified services running on the host. The reporting command is     
  designed to be run as a cron job or systemd timer.                            

  - Send inventory to the server:                                               

      $ inventory send 

  - Send inventory to the server with debug logging:                            

      $ inventory send --log-level debug 

  - Send inventory to the server with a custom config file:                     

      $ inventory send --config-file /path/to/config.yaml 

  - Start the server:                                                           

      $ inventory server 

  - Start the server with debug logging:                                        

      $ inventory server --log-level debug 

  - Start the server with a custom config file:                                 

      $ inventory server --config-file /path/to/config.yaml 

## Options

```
  -c, --config-file string    (default "/home/runner/work/inventory/inventory/inventory.yml")
  -h, --help                 help for inventory
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory completion](inventory_completion.md)	 - Generate the autocompletion script for the specified shell
* [inventory config](inventory_config.md)	 - Create an configuration file for the inventory application
* [inventory search](inventory_search.md)	 - search for services or containers
* [inventory send](inventory_send.md)	 - Send host and container information to the server
* [inventory server](inventory_server.md)	 - starts the RPC and HTTP servers

