#!/usr/bin/env sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"

mkdir -p build/acceptance acceptance/generated tmp/go-cache tmp/go-modcache
rm -rf acceptance/generated/*

for feature in features/*.feature; do
  name=$(basename "$feature" .feature)
  json="build/acceptance/${name}.json"
  tmp/bin/gherkin-parser "$feature" "$json"
  env GOCACHE="$ROOT/tmp/go-cache" GOMODCACHE="$ROOT/tmp/go-modcache" \
    go run ./cmd/acceptance-entrypoint-generator "$json" acceptance/generated
done

env GOCACHE="$ROOT/tmp/go-cache" GOMODCACHE="$ROOT/tmp/go-modcache" \
  go test ./acceptance/generated
