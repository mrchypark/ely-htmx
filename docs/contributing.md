# Contributing and verifying changes

These commands are for a repository checkout. Consumers follow
[getting started](getting-started.md) and do not need the Node toolchain.
Use Go 1.26.x and Node 24 to match [.github/workflows/test.yml](../.github/workflows/test.yml).
Check the actual workflow filename if it changes; CI is authoritative.

## Source map

- `ui.go`, `component.go`: public template/assets/component APIs.
- `templates/`, `assets/`: shared markup, behavior, CSS and embedded assets.
- `assets/input.css`: source styles; `assets/app.css`: committed Tailwind output.
- `sdui/`: strict JSON decode, typed registry, rendering, HTTP actions and store.
- `internal/gallery/`: demonstration fixtures/routes, never production defaults.
- `internal/playground/`, `cmd/sdui-playground/`: authoring server.
- `examples/minimal/`, `examples/app/`, `examples/sdui/`: consumer examples.
- `.agents/skills/ely-htmx/`: reusable agent entrypoint.

Fix library behavior in shared templates/assets, not only a gallery wrapper.
Keep unrelated edits intact. Retain licenses and attribution. Write new docs and
visible example copy in English. Preserve established documentation URLs.

## Standard checks

Run from repository root (the order matches CI):

```sh
go test ./...
go vet ./...
go test -race ./sdui ./examples/sdui ./internal/playground
npm ci
npm test
npm run build:css
git diff --exit-code assets/app.css
```

When intentionally changing styles, inspect and commit the generated CSS, then
rerun the build and ensure it is reproducible. Do not treat a generated CSS diff
as automatically safe. For component/contract changes also run:

```sh
go generate ./sdui
go run ./cmd/sdui-catalogue -check
```

Review generated contract/fixture diffs. Update focused renderer, escape,
interaction and regression tests rather than testing only documentation wording.

## Browser and static Pages checks

```sh
npx playwright install --with-deps chromium
npm run test:browser
node scripts/sdui-static-unit-test.mjs
bash scripts/build-pages.sh dist /elysdui/
node scripts/sdui-static-browser-test.mjs
GALLERY_STATIC=1 node scripts/gallery-browser-test.mjs
```

The browser scripts start their own servers; do not leave competing servers on
their ports. Tests may need OS browser dependencies and permission to bind local
sockets. If blocked, report the failing command and rely only on separately
verified CI results, not an assumed browser pass. `dist/` and `test-results/` are
build/test output; do not include them in a source PR.

The static build compiles the Go renderer to WASM and writes a base-path-aware
Pages bundle. Building/testing it does not deploy. Deployment is the main-branch
workflow after all jobs pass; merging/deploying requires maintainer approval.
See [GitHub Pages](github-pages.md) for hosting and preview limits.

## Documentation and skill checks

Run `go test ./examples/minimal` after changing the starter. Copy its source and
test into a new module using [the consumer commands](getting-started.md) to catch
accidental use of internal packages or checkout-only assumptions. When revising
the skill, have an independent agent use only the entrypoint and linked docs to
build a tiny consumer in a temporary directory. Record gaps and fix observed
failures. Check relative links, exact CLI flags and examples against the code.

In a PR, report the exact revision, checks passed, failures/limits and any
remaining manual/browser validation. Historical [v0.1.0 verification](verification.md)
is not evidence for a new commit. Open a draft PR; do not merge or deploy implicitly.

## Native profile slice

For `sdui/native.go`, `examples/native-sdui` or native app changes, also run:

```sh
go test ./sdui ./examples/native-sdui
go test -race ./sdui ./examples/native-sdui
```

Follow [iOS](../native/ios/README.md) and [Android](../native/android/README.md)
for native build/test commands. `.github/workflows/native.yml` compiles the real
app targets, runs platform model/transport tests, and exercises native profile
validation/save interactions against Go in iOS Simulator and Android Emulator.
It retains screenshots and reports; inspect the exact-head run before claiming
these checks passed. Shared JSON fixtures are
verified against actual Go HTTP responses, then decoded by both native clients.
To intentionally regenerate them after a reviewed contract change:

```sh
UPDATE_NATIVE_FIXTURES=1 go test ./examples/native-sdui -run TestNativeFixtures
```

A v1 wire change requires updating all three implementations/fixtures together.
Report app compilation, model tests and actual simulator/device interaction as
separate stages. Native accessibility, IME and focus behavior need manual
platform validation beyond a compile.
