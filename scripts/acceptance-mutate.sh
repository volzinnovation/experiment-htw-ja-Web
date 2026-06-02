#!/usr/bin/env sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"

level="${1:-soft}"

./scripts/acceptance.sh

mkdir -p build/acceptance-mutation tmp/go-cache tmp/go-modcache

for feature in features/*.feature; do
  name=$(basename "$feature" .feature)
  env GOCACHE="$ROOT/tmp/go-cache" GOMODCACHE="$ROOT/tmp/go-modcache" \
    tmp/bin/gherkin-mutator \
      -feature "$feature" \
      -work-dir "build/acceptance-mutation/$name" \
      -generated-dir acceptance/generated \
      -runner-worker "go run ./cmd/acceptance-runner-adapter" \
      -status-interval 5s \
      -level "$level" \
      -workers 1
done
