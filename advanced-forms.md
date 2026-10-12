# Advanced forms API

Alias components preserve their native control semantics without requiring fixture-only `Type` values. `DateTimePicker`, `TimePicker`, `WeekPicker`, `YearPicker`, and `RadioCard` render their matching native input types. `QuarterPicker` renders Q1–Q4 with an optional `Name_year` numeric field when `Year` or `YearLabel` is present. `TimezoneSelect` renders a native select. `DurationPicker` renders a non-negative numeric value and a unit select from `UnitsOptions`, falling back to `Options`.

`RangeSlider` keeps its start value at or below its end value. The client clamps an end value below the start and updates each native range bound.

`OnboardingWizard` disables inactive steps. It submits values from every visited step in the existing `finish` event payload. An empty wizard shows its empty text and disables its controls.

`ImportExportSettings` validates the current JSON before it emits `import-settings` or `export-settings`. Each action event contains the parsed caller data in `detail.data`; it does not claim that settings were saved.

`ApiKeyManager`, `InvitationList`, `PaymentMethodForm`, and `Paywall` have dedicated markup and emit their caller-provided actions. `WhatsNewDialog` uses a native dialog. Its instance-local trigger opens it, `Open` opens it initially, and its close button emits `dismiss`.
