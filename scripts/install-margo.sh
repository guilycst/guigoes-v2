#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
# Stable release with footnote navigation and mobile viewer fixes.
version=v0.0.28
mkdir -p .tools
if test -x .tools/margo/margo && .tools/margo/margo version | head -1 | grep -qx "margo $version"; then
  exit 0
fi
stage=$(mktemp -d "$PWD/.tools/margo-build.XXXXXX")
trap 'rm -rf "$stage"' EXIT HUP INT TERM
GOBIN="$stage" GOWORK=off CGO_ENABLED=0 go install "github.com/araihu/margo/cmd/margo@$version"
mkdir -p .tools/margo
mv "$stage/margo" .tools/margo/margo
