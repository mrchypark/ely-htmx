# Communication component API

All components use `ui.Component` with `Chapter`, `Name`, `ID`, and flat `Props`. The wrapper supplies the instance ID and the `communication` family. Strings are escaped by Go templates.

Shared props: `Title`, `Label`, `Content`, `Value`, `Name`, `Items`, `Options`, `Disabled`, and `Required`. List items accept `Value`, `Label`, `Title`, and `Content`. Options accept `Value`, `Label`, `Selected`, and `Disabled`.

Chat messages accept `Author`, `Time`, `Attachments`, `Context`, and `Commands`. Dispatch `ely:stream` to a StreamingText root with `{text, status}` to update its local displayed stream. Agent timelines use item `Status`; tool approval keeps the supplied `Payload` in the action value. Artifact options update the local displayed version label only.

Mail uses `To`, `From`, `Subject`, and `Value`. Calendar day objects use `Date`, `Day`, and `Events`; events use `Value`, `Label`, `Start`, and `End`. EventEditor prevents submission when `End` is not after `Start`. Project boards use `Lanes`, each with `Value`, `Label`, `Cards`, and card `Value`, `Title`, `Priority`, and `Due`. Drag or Move emits `move-task` with task and lane IDs. Gantt, workload, milestone, and timer views use caller progress and time values.

Native forms emit their `data-ely-submit` action with full `FormData` and `values`. Buttons emit their declared `data-ely-action`. Selection, ranges, filters, checkboxes, and toggles emit `ely:change`; click actions include selected IDs, dates, task IDs, and lane IDs where relevant. The host must complete provider work. The components only update local display state.
