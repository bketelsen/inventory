#!/bin/sh
set -e

SED="sed"
if which gsed >/dev/null 2>&1; then
	SED="gsed"
fi
NEXT=`svu`
wholething="# inventory <small>$NEXT</small>"
# update this directory to the default value
# of the `--output` flag on the doc generation command
# and at the end of the script too
rm -rf ./docs/inventory*.md
NOCOLOR=1 go run ./cmd/inventory gendocs
"$SED" \
	-i'' \
	-e 's/SEE ALSO/See also/g' \
	-e 's/^## /# /g' \
	-e 's/^### /## /g' \
	-e 's/^#### /### /g' \
	-e 's/^##### /#### /g' \
	./docs/inventory*.md
echo $NEXT
"$SED" \
	-i'' \
	 "/v[0-9]\+\.[0-9]\+\.[0-9]/c $wholething" \
	./docs/_coverpage.md

NOCOLOR=1 go run ./cmd/inventory changelog -n 4
