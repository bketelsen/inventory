# inventory config

Create an configuration file for the inventory application

## Synopsis

Create an configuration file for the inventory application.

This will create a file named inventory.example.yaml in the current directory.
The file will contain the following sections:
- server.address 	
	* the IP:port of the inventory server
- verbose 		
	* true/false - whether to print verbose output
- location 		
	* the location of the server (freeform text)
- description 		
	* the description of the server (freeform text)

In order to use this configuration automatically, you must move it to one of 
the following locations:

- /etc/inventory/
- ~/.config/inventory/

The file must be named "inventory.yml" to be picked up automatically.

Example:
inventory config
sudo mkdir -p /etc/inventory
sudo mv inventory.example.yaml /etc/inventory/inventory.yaml

Be sure to edit the file to set your actual server address and location.
The server.address is the IP:port of the inventory server.

```
inventory config [flags]
```

## Options

```
  -h, --help   help for config
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/home/runner/work/inventory/inventory/inventory.yml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory](inventory.md)	 - Inventory is a tool to collect and report deployment information

