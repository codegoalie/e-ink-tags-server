#!/usr/bin/env sh

set -e

rm etag-server || true
env GOARCH=arm64 GOOS=linux go build -o etag-server
