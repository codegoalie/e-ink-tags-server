#!/usr/bin/env sh

set -e

./build.sh

scp etag-server filch:etag-server/
