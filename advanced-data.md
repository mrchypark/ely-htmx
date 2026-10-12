# Data family API

The data family renders `data`, `lists`, `tables`, `charts`, `finance`, and `dashboard` catalogue stories. Each instance uses its wrapper `ID` and keeps interaction state inside that root.

All stories accept caller-owned values only. List-oriented stories accept `Items`, with `Value`, `Label`, `Detail`, `Selected`, and `Disabled`. Table stories accept `Columns` (`Label`) and `Rows` (`Cells`). Trees use recursive `Items[].Children`.

Charts require caller geometry. Line charts use `Points` and optional `GridPath`. Areas use `AreaPath` and `Bars`. Pie charts use `Segments[].Path`; scatter charts use point `X/Y/R`; radar charts use polygon point strings; heatmaps use `Cells[].Level`; treemaps use positioned items; Sankey charts use `Links[].Path` and positioned `Nodes`; waterfalls and histograms use `Bars`; Gantt charts use `Tasks[].Left/Width`; networks use `Edges` and positioned `Nodes`; bullet charts use `Value`, `Max`, and `Target`. There is no fallback chart series.

RowGrouping accepts Groups. Each group has Value, Label, and canonical table Rows. Each group starts expanded and has its own disclosure control. Callers with a flat Rows[].Group collection must group it before rendering.

TreeTable accepts recursive top-level Rows. Each row has Label, Cells, and optional Children with the same shape. Every node has independent disclosure state. Node IDs derive from the component ID and the row index path, never caller text.

PivotTable accepts Measures. Each measure has Value, Label, Selected, and canonical table Rows. Selecting a measure emits ely:change. Submitting emits the pivot-measure action with the native measure payload. The host can select a measure with ely:data detail {measure}. The component does not calculate aggregates.

Candlestick charts use `Candles` with wick and body geometry. Depth charts use `BidPath` and `AskPath`. Indicator overlays use `Paths`. Drawing tools use caller-owned `Tools` and `Drawings`. Order books use `Bids` and `Asks`. QR and barcode output uses the caller `Value` through the shared encoders.

`PropertyGrid` accepts `Fields` (`Name`, `Label`, `Value`, `Required`, `Disabled`). `Spreadsheet` accepts `Cells`. `OrderEntry` accepts `Symbol`, `Quantity`, and `Price`. `LeverageSlider` accepts `Value` and `Max`. `CurrencyConverter` accepts `From` and `Value`. `RealtimeChart` accepts `Samples`. `Carousel` accepts `ImageURL`, `ImageAlt`, and `Caption`.

Native changes emit `ely:change`. Selection, reordering, filtering, column visibility, range updates, carousel navigation, chart samples, drawing-tool selection, sorting, and alert acknowledgement also emit `ely:change` with action-specific detail. Forms emit host actions through `data-ely-submit`: `save`, `search`, `order`, `convert`, `screen`, and `indicators`. The order form validates symbol, positive quantity, and non-negative price before it asks the host to review the payload. It does not report a completed transaction.

`Tree` and `FileTree` disclose recursive `Items[].Children`. `CheckboxTree` gives each recursive node the caller `Name`; node `Disabled` is honored. `TreeSelect` includes descendants with visible depth. `GroupedList` uses `Groups[].Label/Items`. `InfiniteList` emits `load-more` with caller `Cursor`; the host supplies further rows. `SwipeableListItem` supports a left swipe or accessible Dismiss button, emits the caller `Value`, and provides local Undo. `VirtualList` and `VirtualTable` bound the scrolling viewport and use browser `content-visibility: auto` to defer offscreen rendering; they retain all caller rows in the DOM rather than removing them from a JavaScript window.
