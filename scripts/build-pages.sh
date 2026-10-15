#!/usr/bin/env bash
set -euo pipefail
# Pure static build: no remote APIs, credentials, or runtime server are needed.
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"
out="${1:-dist}"
base="${2:-/ely-htmx/}"
if [ -e "$out" ]; then
  echo "Output directory already exists: $out (choose a fresh directory)" >&2
  exit 1
fi
mkdir -p "$out"
go run ./cmd/ely-gallery-static -dir "$out" -base "$base"
go run ./cmd/sdui-playground-static -dir "$out/playground"
GOOS=js GOARCH=wasm go build -trimpath -ldflags='-s -w' -o "$out/playground/renderer.wasm" ./cmd/sdui-playground-wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$out/playground/wasm_exec.js"
node scripts/version-pages-assets.mjs "$out" "$base"
touch "$out/.nojekyll"
printf 'Static gallery and browser-only editor built in %s\n' "$out"
node -e 'const fs=require("fs"),zlib=require("zlib");const b=fs.readFileSync(process.argv[1]);console.log(JSON.stringify({wasmBytes:b.length,gzip9Bytes:zlib.gzipSync(b,{level:9}).length}));' "$out/playground/renderer.wasm"
