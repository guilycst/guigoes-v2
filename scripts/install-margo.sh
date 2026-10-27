#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
version=0.0.25
os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$(uname -m)" in
  arm64|aarch64) arch=arm64 ;;
  x86_64) arch=amd64 ;;
  *) echo 'Unsupported CPU architecture' >&2; exit 1 ;;
esac
case "$os" in darwin|linux) ;; *) echo 'Unsupported OS' >&2; exit 1 ;; esac
mkdir -p .tools
if test -x .tools/margo/margo && .tools/margo/margo version | head -1 | grep -qx "margo v$version"; then
  exit 0
fi
stage=$(mktemp -d .tools/margo-download.XXXXXX)
trap 'rm -rf "$stage"' EXIT HUP INT TERM
archive="margo_${version}_${os}_${arch}.tar.gz"
base="https://github.com/araihu/margo/releases/download/v$version"
curl -fsSL --retry 3 "$base/$archive" -o "$stage/$archive"
curl -fsSL --retry 3 "$base/checksums.txt" -o "$stage/checksums.txt"
(
  cd "$stage"
  awk -v name="$archive" '$2 == name { print }' checksums.txt > selected.sha256
  test -s selected.sha256
  if command -v sha256sum >/dev/null; then sha256sum -c selected.sha256; else shasum -a 256 -c selected.sha256; fi
  tar -xzf "$archive" margo
)
mkdir -p .tools/margo
mv "$stage/margo" .tools/margo/margo
