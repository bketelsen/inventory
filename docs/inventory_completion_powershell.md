# inventory completion powershell

Generate the autocompletion script for powershell

## Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	inventory completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
inventory completion powershell [flags]
```

## Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/home/runner/work/inventory/inventory/inventory.yml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [inventory completion](inventory_completion.md)	 - Generate the autocompletion script for the specified shell

