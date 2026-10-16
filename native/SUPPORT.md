# Native component coverage

[The machine-readable inventory](support-matrix.json) is an exact allowlist of
implemented semantic node types for SwiftUI and Compose. All unlisted types are
unsupported, including qualified web catalogue component names. The Go resolver
and both clients reject unsupported required components; there is no WebView or
placeholder success fallback. The inventory is checked against all three
capability lists by the Go tests.

Only seven types are currently implemented: region, form, stack, text, input,
button and notice. Each inventory entry identifies its real platform control and
supported modes. Profile v1 remains a single-field, single-action vertical slice.
A shell-level scroll view or loading indicator does not count as support for a
server-driven scroll or progress node. A web catalogue entry does not count as a
native renderer implementation. GPUI coverage is not established by this inventory.

The next coverage increments are multi-field forms with independently reconciled
edit revisions, toggle/checkbox, select/picker and composition-safe multiline
input, followed by lists, images, dividers, progress and adaptive containers.
These are planned, not implemented. Adding a capability requires the Go resolver,
wire contract, both native decoders/renderers, fixtures and tests to agree.

## Verification boundaries

The Native profile workflow runs Go protocol tests, Swift model/transport tests,
iOS app compilation and simulator interaction, plus Kotlin model/transport tests,
Compose app compilation and Android emulator interaction. A green run proves only
its exercised behaviors on that exact commit. It does not certify every property,
manual accessibility, every IME, or physical devices. Consult the PR checks for
current results; this inventory deliberately contains no historical pass badge.
