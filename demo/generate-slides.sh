#!/bin/sh

docker run --rm \
    -v "$PWD:/app:Z" \
    -e LANG="$LANG" \
    -w /app \
    -u $(id -u):$(id -g) \
    --entrypoint marp-cli.js \
    marpteam/marp-cli \
    "--html" "README.md"

