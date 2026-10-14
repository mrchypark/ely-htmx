# SDUI profile example

Run from the repository root with Go 1.26 or newer:

```sh
PORT=8093 go run ./examples/sdui
```

Open <http://127.0.0.1:8093>. The server binds only to loopback. Enter a display
name and save; the advanced switch also saves through the same server action.
Choose an optional file (at most 64 KiB) to exercise multipart `FormData`. The
demo records its name and size, then discards its bytes.

The declarative screen is `screen.json`. To load an external screen without
rebuilding the executable:

```sh
SDUI_SCHEMA=./examples/sdui/screen.json PORT=8093 go run ./examples/sdui
```

The file is strictly decoded and compiled before serving. Changes require a
restart. The loader and action names still have to match the host's Go registry.
There is no live reload or remote schema fetch.

## Important demo limits

Visiting the page automatically creates a random, HTTP-only, same-site local
session. This demonstrates identity isolation and CSRF checks; it is **not
production authentication**. Replace it with real application authentication
and resource-level authorization before using this pattern beyond localhost.

Profile state and request replay records are in memory, bounded, and expire.
They disappear on restart and are not coordinated across processes. The demo
allows at most 256 sessions for one hour, 512 replay entries for ten minutes,
and a 128 KiB HTTP request body. A production mutation needs transactional,
durable idempotency in its own datastore. Cancellation never rolls back a
committed action.

The hidden `delay` field delays responses by at most 300 ms so the browser test
can verify out-of-order delivery. It does not delay the state commit. This is a
local demonstration aid, not an application API feature.

## Checks

```sh
go test ./sdui ./examples/sdui
npm ci
npx playwright install chromium
node scripts/sdui-browser-test.mjs
```

The browser script starts its own Go server on an available loopback port and
closes both server and browser afterward. See [the SDUI guide](../../docs/sdui.md)
for the schema, registry, host API, security boundary and supported scope.
