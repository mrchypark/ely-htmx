# Locate documentation and choose a workflow

This skill is portable: copy the entire `ely-htmx` directory, including this
reference. Its knowledge source is the matching elysdui repository/module, not
an assumed sibling path in your application.

## In the elysdui repository

The root contains `go.mod` with module `github.com/mrchypark/ely-htmx` and
`docs/README.md`. Start there. Read only references needed for the task:

- `docs/getting-started.md`: pinned dependency, standalone consumer commands,
  copyable `examples/minimal`, naming/property surfaces, assets and troubleshooting.
- `docs/agent-sdui.md`: schema, typed bindings, action/loader registration, state,
  security, CLI and playground caveats; links to full contracts.
- `docs/contributing.md`: source ownership, exact CI commands, generated artifacts.
- `docs/asset-loading.md`: family/dependency loading and HTMX fragment constraints.
- `docs/sdui-transport.md`: compression safety and transport settings.
- `docs/sdui-playground.md` and `docs/github-pages.md`: local versus static editor.
- `COMPONENTS.md`, `advanced-*.md`, `themes.md`: detailed component/theme APIs.

## In a consumer application

Inspect `go.mod` and `go list -m -json github.com/mrchypark/ely-htmx` first. The
returned `Dir` locates the selected module source (or local replacement); read
its docs and public APIs there. Do not edit the module cache. If missing, use
`go mod download github.com/mrchypark/ely-htmx` when network access is available.
For an uninitialized app, follow the pinned-version commands in the repository's
getting-started guide. Match documentation to that version.

A skill copied from a newer checkout can mention features absent from the app's
pinned release. Confirm exported APIs and files before using them. The repository
is [mrchypark/elysdui](https://github.com/mrchypark/elysdui); browse the selected tag/commit, not an
unverified latest branch, when local source is unavailable. Report inaccessible
sources rather than inventing APIs.

## Installing this skill

For repository-local discovery, copy this entire directory into the consumer's
`.agents/skills/ely-htmx/`. For personal Codex discovery, copy it into
`${CODEX_HOME:-$HOME/.codex}/skills/ely-htmx/`. Do not overwrite an existing skill
without reviewing its customizations. Agents without automatic skill discovery
can be explicitly instructed to read `SKILL.md`. The skill is guidance, not a
package installer or permission to publish/merge/deploy.
