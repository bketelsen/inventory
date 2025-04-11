# inventory completion bash

Generate the autocompletion script for bash

## Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(inventory completion bash)

To load completions for every new session, execute once:

### Linux:

	inventory completion bash > /etc/bash_completion.d/inventory

### macOS:

	inventory completion bash > $(brew --prefix)/etc/bash_completion.d/inventory

You will need to start a new shell for this setup to take effect.


```
inventory completion bash
```

## Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/home/bjk/inventory/inventory.yml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory completion](inventory_completion.md)	 - Generate the autocompletion script for the specified shell

