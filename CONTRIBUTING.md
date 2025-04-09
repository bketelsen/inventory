# Contributing

By participating in this project, you agree to abide our
[code of conduct](CODE_OF_CONDUCT.md).

## Set up your machine

`inventory` is written in [Go](https://go.dev/).

Prerequisites:

- [Task](https://taskfile.dev/installation)
- [Go 1.24+](https://go.dev/doc/install)

Other things you might need to run some of the tests (they should get
automatically skipped if a needed tool isn't present):

- [gum](https://github.com/charmbracelet/gum?tab=readme-ov-file#installation)



## Building

Clone `inventory` anywhere:

```sh
git clone git@github.com:bketelsen/inventory.git
```

`cd` into the directory and install the dependencies:

```bash
go mod tidy
```

You should then be able to build the binary:

```bash
go build ./cmd/inventory 
./inventory --version
```

## Testing your changes

You can create a branch for your changes and try to build from the source as you go:

```sh
task build
```

When you are satisfied with the changes, we suggest you run:

```sh
task ci
```

Before you commit the changes, we also suggest you run:

```sh
task fmt
```

### A note about Windows

Make sure to enable "Developer Mode" in Settings.

## Creating a commit

Commit messages should be well formatted, and to make that "standardized", we
are using Conventional Commits.

You can follow the documentation on
[their website](https://www.conventionalcommits.org).

The `commit` task will take care of the hard work for you:

```sh
task commit
```

## Submitting a pull request

Push your branch to your `inventory` fork and open a pull request against the main branch.
