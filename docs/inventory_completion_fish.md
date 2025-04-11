# inventory completion fish

Generate the autocompletion script for fish

## Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	inventory completion fish | source

To load completions for every new session, execute once:

	inventory completion fish > ~/.config/fish/completions/inventory.fish

You will need to start a new shell for this setup to take effect.


```
inventory completion fish [flags]
```

## Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/home/runner/work/inventory/inventory/inventory.yml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory completion](inventory_completion.md)	 - Generate the autocompletion script for the specified shell

