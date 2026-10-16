# Working with elysdui

Read [.agents/skills/ely-htmx/SKILL.md](.agents/skills/ely-htmx/SKILL.md) for tasks
that consume or modify this library. Use [docs/README.md](docs/README.md) to find
focused references; do not load the full catalogue by default.

- Consumer work: [docs/getting-started.md](docs/getting-started.md).
- Library changes: [docs/contributing.md](docs/contributing.md).
- JSON screens: [docs/agent-sdui.md](docs/agent-sdui.md).

Keep new documentation and user-facing example copy in English. Preserve existing
document URLs and third-party notices. Do not overwrite unrelated working-tree
changes. Verify APIs against the selected module revision; v0.1.0 does not include
SDUI. Use exact component names and case-sensitive property keys.

Shared component behavior belongs in `templates/` and `assets/`, not gallery
fixtures. Host authentication, authorization, persistence and external services
are application responsibilities. Never turn untrusted input into `template.HTML`.
Run the applicable checks in the contributor guide and state exactly what passed,
failed or could not run. Publishing changes does not authorize merging or deploying.
