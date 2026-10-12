# Native component API

All native components use `ui.Component` with the protected `ID`, `Chapter`, `Kind`, `Variant`, and `ComponentName` values supplied by `RenderComponent`.

Common caller properties are `Title`, `Description`, `Label`, `Name`, `Value`, `Placeholder`, `Action`, `ActionLabel`, `Required`, `Disabled`, `Checked`, `RTL`, `Min`, `Max`, `Options`, and `Items`. `Options` accepts objects with `Value`, `Label`, `Selected`, and `Disabled`; `Items` accepts `Value`, `Label`, and `Content`. Strings are escaped.

Buttons and forms emit `ely:action` through the shared API. Local selections, text/range changes, media controls, crop/mask geometry, model orbit/zoom, and stopwatch state emit `ely:change`. The detail contains the component ID, component path, action, and native form data when applicable. File import/export, generation, capture, print/share, system notifications, model downloads, and host operations remain host actions; the renderer never claims that they completed.

`URL`, `Alt`, `Caption`, and `Volume` configure media. A failed or source-less playback remains an error state. `Vertices` and `Faces` are JSON arrays for `3DModelViewer`; `Rates` is a JSON unit-to-base-rate map for `UnitConverter`; `TimeZone` is an IANA zone for `Clock`; `Token` and `Tokens` configure theme editors. `Items` for media and generation grids may add `URL`, `Alt`, and `Progress`.

`ImageGenerationPanel` emits `generate`; `Calculator` emits `calculate`; `UnitConverter` emits `convert`; editable prompt, subtitle, and theme controls emit their named save actions. Player controls emit `play`; timers emit `toggle` and `reset`; crop, mask, and annotations emit structured geometry actions. The host can cancel an action event or bind an operation with `mount` from `/assets/ely-api.js`.
