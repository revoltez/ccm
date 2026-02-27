#!/bin/bash
set -e

echo "=== Building CCM ==="
go build -o bin/ccm ./cmd/ccm/

echo "=== Running CCM with help ==="
./bin/ccm --help
