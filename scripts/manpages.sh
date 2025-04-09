#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run ./cmd/inventory man | gzip -c -9 >manpages/inventory.1.gz
