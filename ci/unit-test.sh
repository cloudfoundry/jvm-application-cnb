#!/usr/bin/env bash

set -euo pipefail

GOCACHE="$PWD/go-build"

cd jvm-application-buildpack
go test ./...
