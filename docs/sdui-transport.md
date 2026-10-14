# SDUI transport and compression

SDUI uses ordinary same-origin HTTP `fetch`: form input goes to a registered
server action and the response is bounded HTML for one declared region. Screen
JSON is decoded/compiled on the server; it is **not sent on every interaction**.
Keep one registry/store per process, rather than rebuilding them per request.

## Defaults

- Public bundled text assets are precompressed once per process as zstd, Brotli
  and gzip. `Accept-Encoding` quality values choose an acceptable representation;
  equal-quality static variants use the smallest available encoded bytes.
  Older clients get gzip or identity. Browser response decoding is native.
- Representation-specific strong ETags and `Vary: Accept-Encoding` keep caches
  correct. Public assets revalidate (`public, no-cache`); warm requests can return
  304. HEAD has the same representation headers and no body. Identity supports
  ranges; encoded responses ignore Range and return the complete representation.
- Text-only requests of at least 1,024 URL-encoded bytes use native
  `CompressionStream('zstd')` when supported, then gzip. Compression is sent only
  when smaller. Unsupported APIs/formats and local compression failures retain
  native FormData. There is no WASM dependency or automatic network retry.
- Small requests and all file uploads retain native multipart FormData. Repeated
  fields, line endings, cancellation and idempotency survive both paths.
- Lottie is loaded once only when an actual animation mounts, including a later
  region swap. A hidden conditional animation no longer adds its initial download.
- Authenticated HTML stays **uncompressed by default**. It carries CSRF tokens and
  can reflect user input, making blanket response compression vulnerable to
  BREACH-style length oracles. `Cache-Control: no-store, no-transform` tells
  compliant intermediaries not to transform it. Configure upstream compression
  accordingly: these headers cannot prevent an incorrectly configured proxy.

## Request security and host controls

Compressed requests use `application/x-www-form-urlencoded` and the
`X-SDUI-CSRF` header. The framework `_csrf` value is excluded from compressed
bytes; the server validates exactly one header, rejects an accompanying body
`_csrf`, then restores the internal field for action compatibility and canonical
idempotency fingerprints. Native requests continue validating the existing body
field. CSRF tokens unsuitable for a header retain the native request path.
Cookies/authorization remain separate HTTP headers. Standard HTTP/2 HPACK and
HTTP/3 QPACK handling still belongs to the browser/server/proxy; body isolation
is not a claim that all possible side channels are eliminated.

Password forms, including externally associated password controls, automatically
skip request compression. Other application secrets cannot be identified reliably:
set `data-sdui-request-compression="off"` on the form, source or ancestor when
fields mix secrets with attacker-controlled data. Host custom templates can use
this attribute; do not assume hidden inputs are safe because they are hidden.

`Config.MaxBodyBytes` bounds decoded input (default 1 MiB).
`Config.MaxCompressedBodyBytes` independently bounds compressed input (defaults
to MaxBodyBytes). Both accept 1 byte through 64 MiB. gzip/zstd are accepted for
URL-encoded forms only; unknown/stacked encodings and compressed multipart are
415, invalid/truncated streams are 400, size violations are 413. zstd additionally
has an 8 MiB window limit, bounded decoder memory and one decoder worker.
Decoded requests have Content-Encoding removed and Content-Length normalized;
there is no double-decompression fallback. Upstream middleware must not silently
decode while retaining the encoding header.

`Config.CompressResponses = true` is an explicit **security-sensitive opt-in**
for hosts that have independently mitigated response compression side channels.
It enables negotiated zstd/br/gzip on rendered HTML and replay responses. It is
not enabled in the demo. `no-store` alone does not mitigate BREACH. Clients that
forbid every available encoding receive 406 before an action can run.

## Why HTTP rather than gRPC-Web?

Browser gRPC-Web is not native gRPC over unrestricted HTTP/2. It adds a browser
protocol/client, generated stubs and commonly a translating proxy. It does not
make this HTML representation smaller by itself. HTTP/2 or HTTP/3 can already
multiplex these ordinary requests when the host/TLS edge supports them; this
library does not force either protocol. Retaining HTTP also preserves forms,
cache revalidation and the existing bounded HTML region contract.

Prioritize smaller region schemas where application semantics allow it. The demo
replaces its editor/results region intentionally to preserve server validation and
field restoration; narrowing that region blindly would change behavior. Existing
limits remain 1 MiB schema, 1,000 nodes, depth 32 and 2 MiB rendered output.

## Measurement

Against the previous PR head, the demo's initial asset closure had 18 resources
(including module imports), 704,614 uncompressed bytes and an eager 168,394-byte
Lottie dependency. The new closure has 17 resources and approximately 540 kB
uncompressed, including the slightly larger progressive adapter. Measured actual
asset handler bodies total approximately 113 kB gzip, 106 kB Brotli and 124 kB
zstd; equal-quality negotiation chooses the smaller variant per asset. These are
payload bytes, excluding HTTP/TLS headers, not measured network latency. Browser
cache state and negotiated protocol still affect end-to-end results.

The initial HTML changes from 10,245 to 10,188 bytes (the eager script removed);
its response remains identity for security. The normal small demo submission does
not cross the compression threshold. Large text forms benefit; files do not.
Response compression levels trade CPU/memory for bytes, so precomputation is
limited to immutable public assets. Registry initialization and render CPU are
separate costs; this patch makes no render-speed claim.

Tests cover native gzip request bytes, progressive zstd selection, missing APIs,
thresholds, uploads, passwords/opt-out, cancellation/retry, encoded/decoded caps,
malformed streams, CSRF ambiguity, replay across representations, qvalues,
ETags/304/HEAD/ranges, and lazy Lottie lifecycle. Hosted Chromium tests additionally
exercise real browser submission and dynamic animation loading.

## Primary references

- [HTTP content negotiation, RFC 9110 §12.5.3](https://www.rfc-editor.org/rfc/rfc9110.html#section-12.5.3)
- [Content-coded entity tags, RFC 9110 §8.8.3.3](https://www.rfc-editor.org/rfc/rfc9110.html#section-8.8.3.3)
- [WHATWG CompressionStream constructor](https://compression.spec.whatwg.org/#dom-compressionstream): unsupported formats throw; detect them at runtime
- [Browser gRPC-Web features](https://github.com/grpc/grpc-web/blob/master/doc/browser-features.md)
- [Official gRPC-Web quick start](https://grpc.io/docs/platforms/web/quickstart/)
- [BREACH mitigation guidance](https://www.breachattack.com/)
