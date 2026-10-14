# Server-driven UI

The optional `github.com/mrchypark/ely-htmx/sdui` package compiles a small,
versioned JSON screen into the existing Ely Go templates. The server owns
business data, validation and actions. A delegated browser adapter sends native
`FormData` and replaces one explicitly named region with server-rendered HTML.

This is a complete bounded v1 path: strict schema decoding → registry-backed
compilation → host loader → initial page and family assets → authenticated,
CSRF-protected action → validated same-ID HTML region. It does not turn the whole
component catalogue into an unrestricted remote renderer.

## Run the example

From the repository root, with Go 1.26 or newer:

```sh
PORT=8093 go run ./examples/sdui
```

Open <http://127.0.0.1:8093>. Try an empty display name, a successful save, the
advanced switch, and an optional attachment. Both the basic form submit and
advanced `ely:change` event use the same transport. Its loader rejects
unsupported query parameters with HTTP 400. The result displays repeated
form fields and the native submit button's value, as well as file metadata.

The [example schema](../examples/sdui/screen.json) is embedded by default. To
change a screen without rebuilding a compiled executable, point the host at an
external file:

```sh
SDUI_SCHEMA=./examples/sdui/screen.json PORT=8093 go run ./examples/sdui
```

The file is loaded, strictly decoded and compiled at startup. Invalid input
prevents startup. Restart after editing the file; there is no live schema reload
or remote schema fetch. The [example host](../examples/sdui/main.go) composes a
real loader and action with the reusable HTTP runtime.

**The example is local-only.** Its server binds to `127.0.0.1`; the initial GET
automatically issues an HTTP-only, same-site random demo session. That is an
identity-isolation demonstration, not a production login system. Replace it
with real authentication and resource-level authorization before deployment.
State is isolated per demo session and disappears on restart. The example keeps
at most 256 sessions for one hour and 512 replay entries for ten minutes, accepts
128 KiB HTTP bodies, and only records the name and size of files up to 64 KiB.
It does not persist uploaded bytes. Its hidden response-delay field is capped at
300 ms solely to make out-of-order responses testable.

## Schema contract

```json
{
  "version": "1",
  "id": "settings-screen",
  "loader": "load-settings",
  "root": {
    "id": "settings-region",
    "type": "region",
    "props": { "label": "Settings" },
    "children": [
      {
        "id": "settings-form",
        "type": "form",
        "event": { "on": "submit", "action": "save-settings", "target": "settings-region" },
        "children": [
          {
            "id": "settings-name",
            "type": "input",
            "props": { "name": "name", "label": "Display name" },
            "bindings": { "value": "profile.name" }
          },
          {
            "id": "settings-save",
            "type": "button",
            "props": { "type": "submit", "label": "Save" }
          }
        ]
      },
      { "id": "settings-message", "type": "text", "bindings": { "text": "message" } }
    ]
  }
}
```

- `version` must be the string `"1"`
- Every screen and node has a unique, stable ID: a letter followed by letters,
  numbers, underscores or hyphens, at most 128 characters. The
  `sdui-internal-` prefix is reserved for generated control IDs
- `type` is a key in the host's explicit registry. Only structural nodes can
  have `children`; nested forms are rejected
- `props` contains registered, typed scalar properties. Strings, booleans,
  integers and finite numbers are validated without HTML coercion
- `bindings` maps a property to a dotted path through `map[string]any` loader or
  action data. No expressions, functions, reflection, array indexing or code
  evaluation are supported. Missing data and incorrect types fail rendering
- A property cannot have both a literal and binding. A binding suppresses that
  property's default; the host must supply its value
- `event` is optional and singular. Its event type must be supported by that node
  definition. `action` is a host registry key; `target` must be a `region` node
  in the same screen and explicitly allowed for that action
- Loader and action names are capability keys, never URLs

`Decode(io.Reader)` rejects unknown JSON fields, duplicate keys, trailing input
and schemas over 1 MiB. Compilation limits trees to 1,000 nodes and 32 levels;
individual strings are limited to 16 KiB. Rendering has a separate 2 MiB output
budget. These bounds do not replace host-level request and business limits.

The default registry is deliberately small:

| Type | Purpose | Supported event |
| --- | --- | --- |
| `stack` | Layout with enumerated gap, direction and alignment | none |
| `region` | Stable same-ID partial replacement boundary | none |
| `form` | Native form containing controls | `submit` |
| `text` | Escaped paragraph text | none |
| `notice` | Existing alert template, optional title and urgency | none |
| `input` | Existing basic input, including file and hidden types | `change` |
| `button` | Existing basic button, optional name and value | `click` |
| `switch` | Existing advanced `forms/Switch` component | `change`, `ely:change` |

A node has one event, so choose either the native or advanced change binding for
a switch. Binding both forms of an event is not supported. Form submit buttons
normally have no separate click binding. Native successful-control rules apply:
unchecked checkboxes and disabled controls are omitted, repeated names remain
repeated, submitter values are preserved, and files stay `File` entries.
Form field names beginning with `_` are reserved for the transport.

The schema cannot supply raw HTML, arbitrary attributes, class names, CSS,
JavaScript, URLs, `hx-*` configuration, template names or catalogue selectors.
Unregistered props are rejected. Existing template escaping remains in force.
A schema's display strings are data, not trusted markup.

## Compile once and connect a host

```go
registry, err := sdui.DefaultRegistry()
if err != nil { return err }
screen, err := sdui.Decode(schemaReader)
if err != nil { return err }
plan, err := sdui.Compile(screen, registry)
if err != nil { return err }

handler, err := sdui.NewHandler(plan, sdui.Config{
    Endpoint: "/settings",
    Authenticate: authenticateRequest,
    Authorize: authorizeSettingsScreen,
    CSRF: csrfTokenForRequest,
    Loaders: map[string]sdui.Loader{
        "load-settings": loadSettings,
    },
    Actions: map[string]sdui.Action{
        "save-settings": {
            Targets: []string{"settings-region"},
            Authorize: authorizeSettingsMutation,
            Run: saveSettings,
        },
    },
})
if err != nil { return err }
if err := ui.BindAssets(mux); err != nil { return err }
mux.Handle("/settings", handler)
```

The compiled `Plan` snapshots the schema and registry definitions. Do not mutate
the trusted template set after constructing a registry, or mutate request data
while rendering it. A plan can be shared between concurrent requests.

Host callbacks:

```go
type Loader func(context.Context, sdui.Request) (map[string]any, error)

// Authenticate(*http.Request) (principal string, err error)
// Authorize(context.Context, principal string) error
// CSRF(*http.Request) (expectedToken string, err error)
// Action.Authorize(context.Context, sdui.Request) error
// Action.Run(context.Context, sdui.Request) (sdui.Result, error)
```

Authentication, screen authorization, per-action authorization and CSRF token
issuance are mandatory. `CSRF` returns the expected token from trusted session
state; it must not echo untrusted submitted input. Keep the token stable for the
rendered screen; rotating it requires a full-page render. The runtime compares the
submitted token in constant time. Use HTTPS and appropriate production cookie
settings in the host, and authorize the actual resource, not merely the fact
that someone is signed in. Per-action authorization receives parsed values and
files before replay lookup, so it can check any requested record against the
authenticated principal on both first execution and retries. The demo only
edits its own session; it never accepts a user-supplied principal.

`Request` contains the authenticated `Principal`, parsed `Values` (`url.Values`),
`Files` (`map[string][]*multipart.FileHeader`) and the original `HTTPRequest`.
Read or copy retained files inside the action; multipart resources are cleaned
up when the request ends. Validate field multiplicity, file contents and
business constraints in the host. Do not trust file names or content types.

Return `Result{Data: data}` for success, or
`Result{Data: data, Status: http.StatusUnprocessableEntity}` for validation. The
entire target subtree must have the binding data it needs. Validate before
mutating state: returning 422 does not roll back host-side writes. Only status
200 (including zero/default) and 422 are accepted as rendered action outcomes.
Callback/render failures return generic failures rather than partial HTML.

The runtime serves an initial full document on GET and only the requested
region on POST. Assets are bound separately with `ui.BindAssets`. The current
asset URLs are root-mounted under `/assets/`; a subpath deployment needs a
matching host asset route. `Plan.Render`, `RenderRegion` and `RenderAssets` are
also available for hosts that own their page shell.

## Registry extensions

The host can opt in more existing basic templates or catalogue components with
`Registry.Register(Definition)`. A definition selects exactly one trusted basic
template, catalogue component or structural primitive. It declares its scalar
properties, property-to-template mappings, enum tokens or numeric bounds, and
supported event types. The schema only sees the registered alias.

For example, add a safe alias for an existing foundation component:

```go
err := registry.Register(sdui.Definition{
    Name: "host-label",
    Chapter: "typography",
    Component: "Label",
    Properties: map[string]sdui.Property{
        "text": {Type: sdui.String, Target: "Label", Required: true},
    },
})
```

Definitions and templates are trusted application code. A definition should
expose only props its template actually supports. The package rejects dangerous
passthrough keys even in definitions, but does not audit arbitrary custom
renderer logic. Extending a registry does not grant an action: the host must
also allowlist its action and target when constructing the handler.

The asset union is computed for the whole compiled screen. The initial page
loads the needed family entrypoints, their stylesheet dependencies and common
control assets once; fragments carry no scripts. There is no runtime discovery
of new schema families after loading a page. Recompile and load a new page when
the supported screen changes.

## Transport, swaps and lifecycle

`assets/sdui.js` is optional and separate from the existing component APIs. It
uses one delegated transport path for basic events and advanced `ely:*` events.
It deliberately sends vetted `fetch` requests instead of allowing schema-defined
`hx-*` attributes: this keeps native `FormData` intact and checks request order
before any DOM mutation.

The adapter reads only server-rendered `data-sdui-*` bindings. It POSTs to a
same-origin absolute endpoint with cookies, an `Idempotency-Key`, and the
reserved `_action`, `_target`, `_csrf` and `_request` form fields. It overwrites
reserved fields rather than allowing a control to override them. It does not
follow redirects for mutation requests.

- Repeated in-flight submits/clicks/actions from the same source are ignored
- A newer request for the same target aborts the older one when possible
- A latest-request/target-identity guard also discards stale responses even when
  abort cannot stop a response that has already committed
- Successful 200 and validation 422 responses must contain exactly one element
  with the requested ID and `data-sdui-region`; malformed or executable response
  markup is rejected before replacement
- The region is replaced with `outerHTML` semantics. The adapter calls available
  HTMX processing and emits `htmx:after:swap` so Ely components initialize and
  detached component instances are cleaned up
- Focus is restored to a matching control by stable ID or field identity
- Values and checked state inside the replaced region are server-authoritative.
  Focus restoration does not restore stale values. Drafts outside the target
  remain untouched; there is no generic dirty-field merge
- Accepted success emits `sdui:success`; 422 emits `sdui:error` with
  `validation: true`. Other transport failures emit `sdui:error` without swapping

Advanced actions are selected by the node's declared event type and action key,
not arbitrary `event.detail.action`. A custom multi-action widget therefore
needs an explicit host-owned adapter or separately registered nodes. A generic
mapping from arbitrary component action names to server capabilities is not
provided.

## Idempotency and failure boundaries

The HTTP handler requires a safe `Idempotency-Key` for every action. Its replay
key is scoped to principal, action and target. It fingerprints canonical form
values and uploaded contents, so multipart boundary changes do not change the
request identity. Identical concurrent or repeated requests receive the cached
outcome; reusing a key with different input returns 409. Authentication and
authorization are checked again before replay.

The built-in cache is **single-process, in-memory and best-effort**. Defaults are
a ten-minute TTL, 1,024 entries and a 64 MiB replay-byte budget. `Config` exposes
`IdempotencyTTL`, `MaxIdempotencyEntries`, `MaxReplayBytes` and `MaxBodyBytes`
(default request limit 1 MiB). Capacity is checked before running an action;
full capacity returns 503. Completed failures are also cached for their key.
Expired entries are cleaned up on subsequent action requests.

This is not durable or distributed exactly-once execution. Restart, expiry,
multiple instances, or a failed response after a host commit can require a new
request. Production business actions must implement durable idempotency with
the same transaction as the mutation. Cancellation is a network/lifecycle
signal, not a rollback. The browser's stale-response guard protects the DOM,
not the order of business operations; hosts must enforce any required versions
or concurrency checks themselves.

Other common status codes are 401 for missing authentication, 403 for denied
authorization or invalid CSRF, 400 for malformed/unknown action input, 413 for an
oversized body, and 500 for host or rendering failures. Do not use error text as
a machine-readable application protocol.

## Supported scope

V1 supports scalar bindings, structural composition, one declared event per
node, host loaders/actions, initial family assets and same-ID region updates.
It does not provide loops, conditions, slots, arbitrary nested prop objects,
remote templates, a client expression language, schema hot reload, history or
routing, optimistic state reconciliation, tenant-aware durable replay storage,
or a generic multi-action widget bridge. JavaScript is required for the SDUI
action flow; the separate [basic app example](../examples/app) demonstrates
progressive full-page form handling.

## Verification

```sh
go test ./sdui ./examples/sdui
go test -race ./sdui ./examples/sdui
npm ci
npx playwright install chromium
node sdui_runtime_test.mjs
node scripts/sdui-browser-test.mjs
```

The real-browser script launches the Go example on an available loopback port
and uses the actual bundled modules. It checks initial selective assets,
validation and successful saves, advanced change events, repeated region swaps
and cleanup, focus, duplicate submit suppression, files/repeated fields/native
submitter data, authentication/CSRF rejection and idempotent replay. Its race
probe temporarily disables transport cancellation while the real server
returns responses out of order, independently proving the pre-swap sequence
guard. All traffic remains local. `SDUI_SCREENSHOT=/path/to/image.png` optionally
saves a full-page screenshot after the checks. Browser installation or launch
failures are verification blockers, never a successful browser result.

### Regression matrix

| Layer | Cases | Test source |
| --- | --- | --- |
| Decode/compile | Unknown versions/fields/components/props, duplicate JSON keys and IDs, nesting/node/string limits, reserved transport names, unsafe bindings, unsupported events/targets, exact numeric conversion | `sdui/schema_test.go`, `registry_test.go` |
| Render/assets | Escaping and rejected trusted-HTML values, missing/wrong bindings, immutable concurrent rendering, same-ID fragments, output limits, reserved control IDs, all seven family asset dependencies | `sdui/render_test.go` |
| HTTP/security | Authentication and record authorization before replay, CSRF/multiplicity/query rejection, body/upload limits, unknown action/target, repeated and 32 concurrent identical requests, key conflicts and principal isolation, cancellation/panic handling, cache count/byte/TTL bounds | `sdui/http_test.go` |
| Browser adapter | Multipart fields/files/submitter, standalone checked/unchecked advanced controls, repeated submits, independent regions, stale fetch and response-body races, malformed/executable responses, focus, swap delegation, pagehide/teardown | `sdui_runtime_test.mjs` |
| Real browser | Actual Go HTTP server and bundled modules, asset responses, validation/success, advanced switch, repeated swaps/cleanup, files, focus, duplicate suppression, deliberately reordered responses, retries, Back/Forward and subsequent submission | `scripts/sdui-browser-test.mjs` |
| Demo host | External schema startup validation, session isolation/expiry, validation errors, attachments, URL rejection, replay and auth/CSRF negatives | `examples/sdui/main_test.go` |
| Fuzz entrypoint | Decode/compile/render arbitrary input without panic or partial output | `sdui/fuzz_test.go` |

Run a bounded exploratory fuzz pass separately with
`go test ./sdui -run='^$' -fuzz=FuzzDecodeCompileRender -fuzztime=10s -parallel=2`.
The ordinary test suite executes its seed corpus. These checks cover the stated
v1 contracts; they are not proof of every host application or arbitrary custom
template. The repository CI also runs the existing Go and JavaScript suites,
Go race checks for SDUI and its example, CSS rebuild verification, and all three
real-browser regression scripts.
