#!/usr/bin/env bash
set -euo pipefail
# Pure static build: no remote APIs, credentials, or runtime server are needed.
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"
out="${1:-dist}"
base="${2:-/elysdui/}"
source "$root/scripts/pages-base.sh"
if ! valid_pages_base "$base"; then
  echo "Base must be a clean absolute path with a trailing slash" >&2
  exit 1
fi
gallery_base="${base}go-htmx/"
if [ -e "$out" ]; then
  echo "Output directory already exists: $out (choose a fresh directory)" >&2
  exit 1
fi
mkdir -p "$out/go-htmx"
cp site/index.html "$out/index.html"
go run ./cmd/ely-gallery-static -dir "$out/go-htmx" -base "$gallery_base"
go run ./cmd/sdui-playground-static -dir "$out/go-htmx/playground"
GOOS=js GOARCH=wasm go build -trimpath -ldflags='-s -w' -o "$out/go-htmx/playground/renderer.wasm" ./cmd/sdui-playground-wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$out/go-htmx/playground/wasm_exec.js"
node scripts/version-pages-assets.mjs "$out" "$base"
touch "$out/.nojekyll"
printf 'Project hub, static gallery and browser-only editor built in %s\n' "$out"
node -e 'const fs=require("fs"),zlib=require("zlib");const b=fs.readFileSync(process.argv[1]);console.log(JSON.stringify({wasmBytes:b.length,gzip9Bytes:zlib.gzipSync(b,{level:9}).length}));' "$out/go-htmx/playground/renderer.wasm"
