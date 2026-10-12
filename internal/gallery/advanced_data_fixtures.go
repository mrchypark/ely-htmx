package uiexamples

import (
	"fmt"

	"github.com/mrchypark/ely-htmx"
)

// advancedDataCases gives the gallery one independent caller-data instance for every data-family primary.
func advancedDataCases() []AdvancedCase {
	stories := map[string][]string{
		"data":      {"Badge", "Tag", "Avatar", "AvatarGroup", "Statistic", "DescriptionList", "PropertyGrid", "Timeline", "ActivityFeed", "Changelog", "Carousel", "Gallery", "BeforeAfter", "Watermark", "QRCode", "Meter", "Stars", "Comparison"},
		"lists":     {"List", "SelectableList", "SortableList", "SwipeableListItem", "VirtualList", "InfiniteList", "GroupedList", "Tree", "CheckboxTree", "TreeSelect", "Outline", "FileTree", "DirectoryListing"},
		"tables":    {"Table", "DataTable", "ColumnPinning", "RowExpansion", "RowGrouping", "DataGrid", "Spreadsheet", "TreeTable", "PivotTable", "VirtualTable", "HeatmapTable"},
		"charts":    {"LineChart", "AreaChart", "PieChart", "ScatterChart", "RadarChart", "HeatmapChart", "Treemap", "SankeyChart", "WaterfallChart", "Histogram", "GanttChart", "NetworkGraph", "Bullet Chart", "RealtimeChart"},
		"finance":   {"CandlestickChart", "IndicatorOverlay", "DrawingTools", "DepthChart", "MarketProfile", "RenkoChart", "MultiChartLayout", "Ticker", "OrderBook", "TimeAndSales", "OrderEntry", "PositionTable", "LeverageSlider", "OptionChain", "Watchlist", "MarketHeatmap", "Screener", "MarketOverview", "EconomicCalendar", "PortfolioSummary", "DividendTable", "CurrencyConverter"},
		"dashboard": {"Dashboard Grid", "UptimeBar", "AlertList", "EventStream"},
	}
	base := map[string]any{
		"Title": "Caller supplied <record>", "Label": "Caller label", "Value": "42", "Description": "Caller supplied detail", "Caption": "Caller supplied chart caption",
		"Columns": []map[string]any{{"Label": "Name"}, {"Label": "State"}},
		"Rows":    []map[string]any{{"Cells": []string{"Alpha", "Ready"}}, {"Cells": []string{"Beta", "Draft"}}},
		"Items":   []map[string]any{{"Value": "alpha", "Label": "Alpha", "Detail": "First item", "Selected": true}, {"Value": "beta", "Label": "Beta", "Detail": "Second item"}},
		"Fields":  []map[string]any{{"Name": "title", "Label": "Title", "Value": "Caller record", "Required": true}},
		"Cells":   []map[string]any{{"Name": "a1", "Label": "A1", "Value": "A1", "Level": "2"}, {"Name": "b1", "Label": "B1", "Value": "B1", "Level": "4"}}, "Points": "0,90 80,60 160,70 240,25 320,40", "Samples": 2,
		"Symbol": "ACME", "Quantity": 3, "Price": "18.50", "From": "USD", "Max": 8,
	}
	var cases []AdvancedCase
	for chapter, names := range stories {
		for index, name := range names {
			props := make(map[string]any, len(base))
			for key, value := range base {
				props[key] = value
			}
			props["Title"] = fmt.Sprintf("%s %s", chapter, name)
			switch name {
			case "Tree", "FileTree", "CheckboxTree", "TreeSelect":
				props["Name"] = "selected_node"
				props["Items"] = []map[string]any{{"Value": "folder", "Label": "Caller folder", "Open": true, "Children": []map[string]any{{"Value": "file", "Label": "Caller file", "Selected": true}}}}
			case "GroupedList":
				props["Groups"] = []map[string]any{{"Label": "Caller group", "Items": base["Items"]}}
			case "InfiniteList":
				props["Cursor"] = "caller-next-page"

			case "LineChart", "RealtimeChart":
				props["GridPath"], props["Points"], props["ViewBox"] = "M0 100H320", "0,90 80,60 160,75 240,30 320,45", "0 0 320 120"
			case "AreaChart":
				props["AreaPath"], props["Bars"] = "M0 100L0 80 100 50 200 70 320 20V100Z", []map[string]any{{"X": 30, "Y": 65, "Width": 24, "Height": 35}}
			case "PieChart":
				props["Segments"] = []map[string]any{{"Label": "Alpha", "Path": "M60 60L60 5A55 55 0 0 1 112 76Z"}, {"Label": "Beta", "Path": "M60 60L112 76A55 55 0 0 1 18 88Z"}}
			case "ScatterChart":
				props["Points"] = []map[string]any{{"X": 30, "Y": 80, "R": 5, "Label": "Alpha"}, {"X": 180, "Y": 35, "R": 8, "Label": "Beta"}}
			case "RadarChart":
				props["GridPoints"], props["Points"] = "100,5 190,70 155,145 45,145 10,70", "100,28 165,76 140,120 60,115 35,72"
			case "Treemap":
				props["Items"] = []map[string]any{{"Label": "Alpha", "X": 0, "Y": 0, "Width": 180, "Height": 120, "TextX": 8, "TextY": 20}, {"Label": "Beta", "X": 182, "Y": 0, "Width": 138, "Height": 120, "TextX": 190, "TextY": 20}}
			case "SankeyChart":
				props["Links"], props["Nodes"] = []map[string]any{{"Path": "M20 30C100 30 140 45 220 45V75C140 75 100 90 20 90Z"}}, []map[string]any{{"Label": "In", "X": 0, "Y": 25, "Width": 20, "Height": 70, "LabelY": 18}}
			case "WaterfallChart", "Histogram":
				props["Bars"] = []map[string]any{{"X": 20, "Y": 50, "Width": 30, "Height": 50}, {"X": 70, "Y": 25, "Width": 30, "Height": 75}}
			case "NetworkGraph":
				props["Edges"], props["Nodes"] = []map[string]any{{"X1": 30, "Y1": 60, "X2": 180, "Y2": 35}}, []map[string]any{{"X": 30, "Y": 60, "R": 10, "Label": "A", "LabelX": 45, "LabelY": 60}}
			case "CandlestickChart":
				props["Candles"], props["Types"], props["Ranges"], props["ViewBox"] = []map[string]any{{"X": 40, "High": 20, "Low": 95, "BodyX": 32, "Open": 45, "Width": 16, "BodyHeight": 30, "Direction": "up"}}, []map[string]any{{"Value": "candle", "Label": "Candle"}}, []map[string]any{{"Value": "1d", "Label": "1 day"}}, "0 0 320 120"
			case "DepthChart":
				props["BidPath"], props["AskPath"], props["ViewBox"] = "M0 100L160 40V100Z", "M160 40L320 100V100H160Z", "0 0 320 120"
			case "OrderBook":
				props["Bids"], props["Asks"] = []map[string]any{{"Cells": []string{"10.00", "12"}}}, []map[string]any{{"Cells": []string{"10.10", "8"}}}
			case "DrawingTools":
				props["Tools"], props["Drawings"], props["ViewBox"] = []map[string]any{{"Value": "line", "Label": "Line", "Selected": true}, {"Value": "erase", "Label": "Erase", "Selected": false}}, []map[string]any{{"Path": "M30 90L250 20"}}, "0 0 320 120"
			case "GanttChart":
				props["Tasks"] = []map[string]any{{"Label": "Research", "Left": 10, "Width": 45, "Detail": "Review"}}
			case "IndicatorOverlay":
				props["Paths"], props["Indicators"] = []map[string]any{{"Path": "M0 80L160 45L320 55"}}, []map[string]any{{"Value": "sma", "Label": "SMA"}}
			case "MultiChartLayout":
				props["Charts"] = []map[string]any{{"Label": "ACME", "ViewBox": "0 0 320 120", "Points": "0,90 160,45 320,60"}}
			case "Screener":
				props["Filters"] = []map[string]any{{"Name": "market", "Label": "Market", "Options": []map[string]any{{"Value": "all", "Label": "All"}}}}
			case "Dashboard Grid":
				props["Widgets"] = []map[string]any{{"Label": "Requests", "Value": "42", "ViewBox": "0 0 320 120", "Points": "0,90 160,45 320,60"}}
			case "UptimeBar":
				props["Services"] = []map[string]any{{"Label": "API", "Value": "99.9%", "Samples": []map[string]any{{"State": "ok"}}}}
			case "Stars":
				props["Stars"], props["AriaLabel"] = []map[string]any{{"Filled": true}, {"Filled": true}, {"Filled": true}, {"Filled": false}, {"Filled": false}}, "3 out of 5"
			case "RowGrouping":
				props["Groups"] = []map[string]any{{"Value": "open", "Label": "Open", "Rows": []map[string]any{{"Cells": []string{"Alpha", "Open"}}}}, {"Value": "done", "Label": "Done", "Rows": []map[string]any{{"Cells": []string{"Beta", "Done"}}}}}
			case "TreeTable":
				props["Rows"] = []map[string]any{{"Label": "Parent", "Cells": []string{"Parent", "Open"}, "Children": []map[string]any{{"Label": "Child", "Cells": []string{"Child", "Ready"}, "Children": []map[string]any{{"Label": "Grandchild", "Cells": []string{"Grandchild", "Draft"}}}}}}}
			case "PivotTable":
				props["Measures"] = []map[string]any{{"Value": "count", "Label": "Count", "Selected": true, "Rows": []map[string]any{{"Cells": []string{"Alpha", "2"}}}}, {"Value": "sum", "Label": "Sum", "Rows": []map[string]any{{"Cells": []string{"Alpha", "9"}}}}}
			case "BeforeAfter":
				props["BeforeURL"], props["AfterURL"], props["BeforeAlt"], props["AfterAlt"], props["Value"] = "/assets/sample-before.svg", "/assets/sample.svg", "Before correction", "After correction", 50
			}
			cases = append(cases, AdvancedCase{Family: "data", Component: ui.Component{Chapter: chapter, Name: name, ID: fmt.Sprintf("audit-%s-%d", chapter, index+1), Props: props}})
		}
	}
	return cases
}
