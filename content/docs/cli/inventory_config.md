---
date: 2025-03-31T20:22:34-04:00
title: "inventory config"
slug: inventory_config
url: /docs/cli/inventory_config/
---
## inventory config

Create an example configuration file for the inventory client (reporter)

### Synopsis

Create an example configuration file for the inventory client (reporter).

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
- ~/.inventory/

The file must be named "inventory.yaml".

Example:
inventory config
sudo mkdir -p /etc/inventory
sudo mv inventory.example.yaml /etc/inventory/inventory.yaml

Be sure to edit the file to set your actual server address and location.
The server.address is the IP:port of the inventory server.

```
inventory config [flags]
```

### Options

```
  -h, --help   help for config
```

### Options inherited from parent commands

```
      --config string   config file (default is /etc/inventory/inventory.yaml)
  -v, --verbose         verbose logging
```

### SEE ALSO

* [inventory](/inventory/docs/cli/inventory/)	 - Inventory is a tool to collect and report deployment information

###### Auto generated by toolbox on 31-Mar-2025
