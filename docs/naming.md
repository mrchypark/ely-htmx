# elysdui naming and compatibility

The product name is **elysdui**, all lowercase. It replaces Ely HTMX in current
project documentation, gallery/playground branding and example application copy.
The name describes a server-driven UI project with a Go/HTML/htmx renderer and a
bounded SwiftUI/Jetpack Compose profile example. It does not imply full native
catalogue coverage or an implemented Ely GPUI runtime adapter.

## Stable integration names

The repository and public website have moved. Package/API identity is retained:

- Repository: `github.com/mrchypark/elysdui`.
- Go module/import path: `github.com/mrchypark/ely-htmx` (unchanged).
- GitHub Pages root: `/elysdui/`, the shared architecture and project hub.
- Gallery, its assets and playground: `/elysdui/go-htmx/`. Old Pages URLs are not guaranteed to redirect.
- CLI: `ely`; existing commands, Go package names and imports are unchanged.
- Reusable agent skill identifier and directory: `ely-htmx`. Its title and
  description use elysdui; retaining the identifier avoids breaking installations
  and existing links to `.agents/skills/ely-htmx/`.
- Swift packages/targets, Android package/application IDs, protocol identifiers,
  schema IDs, CSS classes, JavaScript hooks and asset paths are unchanged.
- Published versions, pinned revisions, historical verification and copyright
  notices remain as originally recorded. Branding does not change v0.1.0 contents.

A future Go module migration needs a separate release/version and downstream
compatibility review. Do not replace module imports or schema identifiers merely
to match the repository name. GitHub repository redirects do not imply Pages
redirects.

## Upstream attribution

[Ely GPUI Component](https://git.zacharyzhang.com/ZacharyZhang-NY/Ely-GPUI-Components)
is the upstream visual/component reference, not a former name for this project.
Keep its name, source-reference fixtures and MIT notices intact. See
[third-party UI notices](../THIRD_PARTY_UI.md). Synthetic catalogue examples may
also retain their own sample names; they are not the elysdui product wordmark.
