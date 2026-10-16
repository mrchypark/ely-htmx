# Ely native iOS profile

This is an actual iOS application target with a SwiftUI screen renderer, not a
WebView. It consumes the shared [profile protocol](../PROTOCOL.md) and the Go
example's resolved seven-component tree. Minimum iOS version: 17.0. The Swift
package requires Swift 5.9 or newer; Xcode 15 or newer supplies the iOS SDK.

## Run on an iOS Simulator

Start the real Go example from the repository root:

```sh
go run ./examples/native-sdui
```

On a Mac with Xcode selected through `xcode-select`, install the standard
[XcodeGen](https://github.com/yonaskolb/XcodeGen) generator, then create the app
project from the checked-in specification:

```sh
brew install xcodegen
cd native/ios
xcodegen generate
open ElyProfile.xcodeproj
```

Choose the `ElyProfile` scheme and an iPhone or iPad simulator, then Run. The
Debug app uses `http://127.0.0.1:8094`; the simulator shares the Mac's loopback
network. Edit the display name, save an empty/one-character name to see native
validation, then save a valid name. Successful saves trim surrounding whitespace.

The app's Reload button asks before discarding unsaved edits. A network failure
retains the exact pending request for **Retry same request**. It does not silently
submit later edits or invent a new idempotency key. You can continue editing
while a save runs; editing is disabled during an explicit reload so fresh edits
cannot be silently overwritten by the load result.

The generated `.xcodeproj` is intentionally ignored: `project.yml` is the source
of truth and creates both the launchable app and its UI-test target reproducibly.
No developer signing identity is needed to build/run the simulator app.

## API origin and local HTTP policy

- Debug's `ELY_API_BASE_URL` build setting defaults to the loopback URL above.
  A Debug launch environment variable with the same name can override it.
- Release has no default origin. Set `ELY_API_BASE_URL=https://your-api.example`
  as an Xcode build setting when integrating a real application-owned backend.
  A missing origin produces a native configuration screen.
- Only Debug's Info.plist contains local-network ATS exceptions. Neither plist
  contains `NSAllowsArbitraryLoads`. Debug declares exact exceptions for
  `localhost`, `127.0.0.1`, and `::1`, including the IP entries required by
  [iOS 17 ATS](https://developer.apple.com/documentation/bundleresources/information-property-list/nsapptransportsecurity/nsexceptiondomains).
  The transport independently rejects HTTP
  unless explicitly enabled for `localhost`, `127.0.0.1`, or `::1`; Release does
  not enable it. Redirects are refused.
- A physical iOS device needs a real HTTPS backend. Do not expose the Go demo's
  automatic authentication/session implementation on a public network.

## Verification

From the repository root, run the Foundation protocol/state/transport tests:

```sh
swift test --package-path native/ios
```

These tests read `native/fixtures/profile-{ready,invalid,saved}.json`, the same
real-handler fixtures used by Go and Android. Keep the repository checkout
intact; the test files resolve the shared fixture directory relative to their
source paths. No fixture data is embedded in the shipping app.

With the Go example running, also exercise real HTTP, cookie retention, native
validation, canonical saving, exact-request replay, and explicit reload:

```sh
ELY_NATIVE_TEST_BASE_URL=http://127.0.0.1:8094 swift test --package-path native/ios
```

The live test is skipped unless that variable is supplied. It uses an isolated
in-memory cookie store. A normal unit-test run also tests streaming rejection of
a response larger than 1 MiB, without needing a server.

Build the actual iOS target, using Xcode and an installed simulator SDK:

```sh
(cd native/ios && xcodegen generate)
xcodebuild -project native/ios/ElyProfile.xcodeproj \
  -scheme ElyProfile -configuration Debug -sdk iphonesimulator \
  -destination 'generic/platform=iOS Simulator' \
  CODE_SIGNING_ALLOWED=NO build
```

Run the actual native UI test against the running Go server. Replace the example
simulator name with an installed device from `xcrun simctl list devices available`:

```sh
xcodebuild -project native/ios/ElyProfile.xcodeproj \
  -scheme ElyProfile -configuration Debug \
  -destination 'platform=iOS Simulator,name=iPhone 16' \
  CODE_SIGNING_ALLOWED=NO test
```

`ElyProfileUITests` verifies native load, invalid save and visible error, clearing
an old error after editing, successful save, canceling reload without losing a
draft, and confirmed reload restoring the server value. It asserts that the app
has no WebViews. A generic simulator build alone is not UI/device verification.

Manual device/simulator checks complement automation:

1. Turn on VoiceOver and large Dynamic Type; labels, inline errors, loading and
   status announcements should remain understandable.
2. Start a Japanese or Chinese IME composition. No request should be sent and
   Save/Retry should remain unavailable until composition commits.
3. Delay a save response while typing a newer edit. The newer edit and focused
   native control should survive the full-region replacement.
4. Stop the Go server during a save. Retry should reuse the pending request;
   restarting this in-memory server invalidates the old session and requires an
   explicit reload rather than silently discarding edits.

## Implementation boundaries

- `ElyNativeCore` contains Codable models, pre-render validation, a value-type
  reconciliation state machine, and Foundation URLSession transport.
- `ElyNativeUI` maps `region`, `form`, `stack`, `text`, `input`, `button`, and
  `notice` to native views. Structural containers use SwiftUI stacks; an input
  uses a native UITextField through `UIViewRepresentable` to observe marked-text
  composition reliably. Buttons, labels, notices and navigation are SwiftUI.
- Stable server IDs identify views; neither dataVersion nor screenRevision is a
  view identity. A focused field is retained across normal region replacements.
  Native input values bind to local drafts, not to the tree's stale `value` prop.
- Validation occurs before tree replacement: 1 MiB decoded response, 1,000 nodes,
  depth 32, unique stable IDs, safe integer revisions, exact registered action
  and field, supported protocol/capabilities, and property types.
- Only a matching pending request/session/screen can update state. Old versions
  are ignored. Canonical values apply only to acknowledged, still-current local
  edit revisions. Conflict/validation preserve drafts; stale errors are hidden.
- URLSession's cookie storage retains/resends the HttpOnly demo cookie across
  requests; the OS handles persistence according to cookie lifetime. Session
  cookies and local drafts are not a durable offline data store. CSRF is sent in
  its own header. Neither sessionID nor CSRF replaces cookie authentication.
- There is no general evaluator, HTML fallback, remote URL action, remote image,
  switch component, catalogue-wide renderer, production auth, or offline queue.
  Unsupported/malformed responses show a native error and explicit reload path
  while retaining existing drafts.

Foundation tests can run on Linux. Linux does not compile the iOS-only renderer
or prove simulator/device behavior; use the Xcode build and UI-test steps above
for those checks.
