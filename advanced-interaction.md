# Interaction API

Render an interaction component with `ui.Component{Chapter, Name, ID, Props}`. `ID` is required, starts with a letter, and is the instance boundary for behavior and events. `Name` accepts a catalogue alias; the renderer receives its primary story as `Kind` and the requested alias as `Variant` and `ComponentName`.

All text props are escaped. `Items` is caller data: each item may have `Value`, `Label`, `Title`, `Content`, `Selected`, and `Href`. Every story accepts `Title`, `Label`, `Content`, `Value`, `Name`, `ActionLabel`, `Action`, `Disabled`, and `Items` where applicable. `Disabled` blocks local behavior and generic host actions. `Required` is used on rendered native form controls. The event root is the parent section, with `data-ely-family="interaction"` and `data-ely-component="<chapter>/<alias>"`; its child renderer retains both primary `data-interaction-kind` and requested alias `data-interaction-variant`.

`ely:action` is a bubbling, cancelable event. Base handling sends `id`, `component`, `action`, button `value`, and native `formData` plus `values` for a `data-ely-submit` form. The family module emits `ely:change` for its local selections and state. Host provider work belongs in `ui.mount`; local controls never claim provider completion.

| Chapter / primary | Props used beyond the common props | Host actions |
| --- | --- | --- |
| buttons/Button | `ButtonVariant`, `Action` | `Action` or `press` |
| buttons/Sizes | `Items[].Size` | `press` |
| buttons/With icons | `Icon` | `press` |
| buttons/Icon button | `Icon` | `press` |
| buttons/ShortcutHint | `Value` | none |
| buttons/ButtonGroup | `Items` | local `ely:change` |
| buttons/SplitButton | `Items` | `Action`, item `Action` or `select` |
| buttons/ToggleButton | `Pressed` | local `ely:change` |
| buttons/ToggleGroup | `Items[].Selected` | local `ely:change` |
| buttons/SegmentedControl | `Items[].Selected` | local `ely:change` |
| buttons/FloatingActionButton | `Icon` | `Action` or `create` |
| buttons/ConfirmButton | `Action`, `Value` | `Action` or `confirm` after hold |
| buttons/CopyButton | `Value`, `Action` | `copy` or `Action` |
| buttons/ActionBar | `Items` | item `Action` or `action` |
| buttons/BulkActionBar | `Items`, `Value` | item `Action` or `bulk` |
| buttons/ActionSheet | `Items` | item `Action` or `select` |
| buttons/QuickActions | `Items` | item `Action` or `action` |
| navigation/Tabs | `Items[].Content`, `Items[].Selected` | local `ely:change` |
| navigation/EditorTabs | `Items[].Selected` | local `ely:change` |
| navigation/Breadcrumb | `Items[].Href` | browser navigation |
| navigation/Pagination | `Items[].Selected` | local `ely:change` |
| navigation/Steps | `Items[].Content`, `Items[].Selected` | none |
| navigation/Sidebar Navigation | `Items[].Href`, `Items[].Selected` | browser navigation |
| navigation/NavigationMenu | `Items[].Href` | browser navigation |
| navigation/BackForwardNavigation | `BackLabel`, `ForwardLabel` | `back`, `forward` |
| navigation/Anchor | `Items[].Href` | browser navigation |
| navigation/JumpTo | `Name`, `Value`, `Required` | `jump` |
| navigation/CommandPalette | `Name`, `Value`, `Items` | item `Action` or `select` |
| navigation/QuickOpen | `Name`, `Value`, `Placeholder`, `Required` | `open` |
| navigation/QuickSwitcher | `Items` | keyboard focus only |
| navigation/SearchPalette | `Name`, `Value`, `Placeholder` | `search` |
| menus/Menu | `Items` | item `Action` or `select` |
| menus/CheckboxMenuItem | `Name`, `Items[].Selected` | native change |
| menus/DropdownMenu | `Items` | `Action`, item `Action` or `select` |
| menus/ContextMenu | `Content`, `Items` | item `Action` or `select` |
| menus/MenuBar Menu | `Items[].Content` | menu open state |
| menus/MegaMenu | `Items[].Content` | none |
| menus/SearchableMenu | `Value`, `Placeholder`, `Items` | local `ely:change` |
| menus/PieMenu | `Items` | item `Action` or `select` |
| overlays/Popover | `Content` | local open and close |
| overlays/HoverCard | `Content` | focus and hover display |
| overlays/Dialog | `Name`, `Value`, `Required`, `Action` | `Action` or `submit` |
| overlays/Lightbox | `Src`, `Alt`, `Content` | local open and close |
| overlays/Spotlight | `Content` | local dismissal |
| overlays/FloatingToolbar | `Items[].Selected` | local `ely:change` |
| overlays/Callout | `Content` | local dismissal |
| overlays/Tooltip | `Content` | focus and hover display |
| feedback/Toast | `Value` | local dismissal only |
| feedback/Notification | `Items` | local dismissal only |
| feedback/Alert | `Urgent`, `Content` | none |
| feedback/Banner | `Content` | none |
| feedback/Callout | `Content` | none |
| feedback/InlineMessage | `Value`, `Content` | none |
| feedback/EmptyState | `Description`, `Content` | none |
| feedback/ErrorView | `Content`, `Action` | `Action` or `retry` |
| feedback/ErrorBoundary | `Content`, `Action` | `Action` or `retry` |
| feedback/Result | `Value`, `Content` | none |
| feedback/ConfirmationCard | `Description`, `Content` | none |
| feedback/SavingIndicator | `Value` | none |
| feedback/Countdown | numeric `Value` | none |
| motion/Spinner | `Label` | none |
| motion/Skeleton | `Label` | none |
| motion/Shimmer | `Content`, `Value` | none |
| motion/LoadingButton | `Loading` | local busy state |
| motion/LoadMore | `Value`, `Action` | `Action` or `load-more`; local count |
| motion/ProgressBar | numeric `Value`, `Max` | none |
| motion/LoadingOverlay | `Content` | none |
| motion/LazyLoad | `Content` | none |
| motion/Suspense | `Loading`, `Content` | none |
| motion/RefreshIndicator | `Action` | `Action` or `refresh` |
| motion/UploadList | `Items` | none |
| motion/Transition | `Content` | none |
| motion/Stagger | `Items` | none |
| motion/AnimatePresence | `Content` | none |
| motion/Flip | `Content` | none |
| motion/Reorder | `Items` | `drag-start`, `drop`; local reorder |
| motion/Marquee | `Content`, `Value` | none |
| motion/Ripple | `Label` | local press effect |
| motion/Glow | `Content` | none |
| motion/Blink | `Content` | none |
| motion/ParticleBackground | `Content` | none |
| interaction/Draggable | `Content`, `DropLabel` | `drag-start`, `drop` |
| interaction/Resizable | `Content` | local `ely:change` |
| interaction/LongPress | `Action`, `Value` | `Action` or `long-press` after hold |
| interaction/Hotkeys | `Action`, `Value` | `Action` or `hotkey` with `key` |
| interaction/Undo | `UndoLabel`, `RedoLabel` | local `ely:change` |
| interaction/VirtualScroller | `Items` | browser scroll |
| interaction/Autofocus | `Name`, `Value`, `Placeholder`, `Autofocus`, `Required` | native change and roving focus |

The gallery fixture supplies every primary with two caller items, a unique `audit-<chapter>-<index>` ID, realistic names and values, and disabled or required cases. Gallery checks verify mounting plus meaningful state changes for tabs, menus, dialogs, selection controls, filters, hotkeys, resize, undo, and native form submission.

`ParticleBackground` and `AnimatedGradient` use local CSS motion. `Lottie Player` uses the bundled `lottie-light` renderer. Supply either `Animation` as inline Lottie JSON or `AnimationURL` as a same-origin JSON URL; `Autoplay` and `Loop` control playback. The renderer provides local play, pause, and seek controls, emits `ely:change` with `name: "frame"`, and destroys the animation when its component unmounts. An unavailable library, invalid JSON, failed request, or cross-origin URL dispatches `ely:error`; it never reports host operation success.
