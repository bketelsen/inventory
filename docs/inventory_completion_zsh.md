# inventory completion zsh

Generate the autocompletion script for zsh

## Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(inventory completion zsh)

To load completions for every new session, execute once:

### Linux:

	inventory completion zsh > "${fpath[1]}/_inventory"

### macOS:

	inventory completion zsh > $(brew --prefix)/share/zsh/site-functions/_inventory

You will need to start a new shell for this setup to take effect.


```
inventory completion zsh [flags]
```

## Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/home/bjk/inventory/inventory.yml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory completion](inventory_completion.md)	 - Generate the autocompletion script for the specified shell

