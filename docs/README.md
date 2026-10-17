# Documentation map

Start with [Using elysdui](getting-started.md) to build an application, or
[Contributing](contributing.md) to change this repository. AI coding agents should
read [AGENTS.md](../AGENTS.md) and the [elysdui skill](../.agents/skills/ely-htmx/SKILL.md).

## Choose a workflow

- **Go templates + htmx:** [minimal integration](getting-started.md), [runnable form](../examples/app), [component API](../COMPONENTS.md).
- **JSON-defined screens:** [agent SDUI workflow](agent-sdui.md), [complete SDUI contract](sdui.md), [composition](sdui-composition.md), [typed catalogue](sdui-catalogue-contracts.md).
- **Native mobile:** [SwiftUI + Compose profile slice](../native/README.md), [supported types and limits](../native/SUPPORT.md), [versioned native protocol](../native/PROTOCOL.md), [multi-field/switch extension](../native/PROTOCOL.md#multi-field-and-switch-extension), [run the controls demo](../examples/native-sdui/README.md#expanded-native-controls).
- **Visual authoring:** [playground](sdui-playground.md), [static GitHub Pages limits](github-pages.md).
- **Assets and performance:** [family loading](asset-loading.md), [transport/compression](sdui-transport.md), [measured benchmarks](sdui-transport-benchmarks.md).
- **Themes:** [theme API](../themes.md).
- **Coverage:** [SDUI implemented/limited behavior](sdui-roadmap.md), [v0.1.0 historical verification](verification.md).

## Component references

The existing detailed references remain at their original paths (some are in
Korean). New onboarding and agent instructions are English. Use exact API names
and typed contracts rather than translating identifiers.

- [Foundation](../advanced-foundation.md)
- [Interaction](../advanced-interaction.md)
- [Data](../advanced-data.md)
- [Workspace](../advanced-workspace.md)
- [Communication](../advanced-communication.md)
- [Native/media](../advanced-native.md)
- [Forms](../advanced-forms.md)
- [Native tooling](../native-tooling.md)

See [component discovery](getting-started.md#discover-components-and-props) before
choosing an alias or assuming a component integrates an external service.

- [Naming and compatibility](naming.md): elysdui branding and stable integration names.

- [Host-owned raster map tiles](map-tiles.md): trusted provider IDs, local fixture, resource limits, and source correspondence.

- [Native maps and canvas design](native-map-canvas-design.md): proposed geographic editing architecture, SDK boundaries, compatibility and acceptance plan; not implemented.
