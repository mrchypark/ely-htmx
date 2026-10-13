package uiexamples

import "testing"

func TestFinanceAndDashboardFixturesKeepSourceSpecimens(t *testing.T) {
	fixtures := map[string]AdvancedCase{}
	for _, fixture := range advancedDataCases() {
		fixtures[fixture.Component.Chapter+"/"+fixture.Component.Name] = fixture
	}
	board := fixtures["dashboard/Dashboard Grid"].Component.Props["Widgets"].([]map[string]any)
	if len(board) != 5 || board[2]["Label"] != "CPU" || board[3]["W"] != 8 {
		t.Fatal("Dashboard Grid lost the five-tile source layout")
	}
	uptime := fixtures["dashboard/UptimeBar"].Component.Props["HealthChecks"].([]map[string]any)
	if len(uptime) != 5 || uptime[4]["Label"] != "Search index" {
		t.Fatal("UptimeBar lost the pending source check")
	}
	incident := fixtures["dashboard/AlertList"].Component.Props["Incident"].(map[string]any)
	if len(incident["Updates"].([]map[string]any)) != 3 {
		t.Fatal("AlertList lost the full source incident timeline")
	}
	orders := fixtures["finance/OrderEntry"].Component.Props
	if orders["Symbol"] != "ELY" || orders["Quantity"] != 100 || len(orders["OrderKinds"].([]map[string]any)) != 3 {
		t.Fatal("OrderEntry lost its caller ticket controls")
	}
	book := fixtures["finance/OrderBook"].Component.Props
	if len(book["Bids"].([]map[string]any)) != 12 || len(book["Asks"].([]map[string]any)) != 12 || len(book["Level2"].([]map[string]any)) != 6 {
		t.Fatal("OrderBook lost its depth and venue context")
	}
	watchlist := fixtures["finance/Watchlist"].Component.Props
	if len(watchlist["Items"].([]map[string]any)) != 8 || len(watchlist["SymbolBadges"].([]map[string]any)) != 4 {
		t.Fatal("Watchlist lost its source symbols or search controls")
	}
	tape := fixtures["finance/Ticker"].Component.Props["Tape"].([]map[string]any)
	if len(tape) != 8 || tape[0]["Value"] != "187.24" || tape[7]["Label"] != "TERRA" {
		t.Fatal("Ticker lost the complete upstream tape")
	}
	portfolio := fixtures["finance/PortfolioSummary"].Component.Props["Holdings"].([]map[string]any)
	if len(portfolio) != 6 || portfolio[0]["Text"] != "$48,210" || portfolio[5]["Label"] != "Cash" {
		t.Fatal("PortfolioSummary lost the caller allocation specimen")
	}
	converter := fixtures["finance/CurrencyConverter"].Component.Props
	if converter["Value"] != 1000 || converter["Converted"] != "€917.43" || len(converter["Transfers"].([]map[string]any)) != 5 {
		t.Fatal("CurrencyConverter lost the source conversion or five transfer specimen")
	}
	leverage := fixtures["finance/LeverageSlider"].Component.Props
	if leverage["MarginUsed"] != "$21,500" || leverage["Risk"] != 13 {
		t.Fatal("LeverageSlider lost margin or bounded risk context")
	}
	calendar := fixtures["finance/EconomicCalendar"].Component.Props
	if len(calendar["Earnings"].([]map[string]any)) != 4 || len(calendar["Stories"].([]map[string]any)) != 4 {
		t.Fatal("EconomicCalendar lost source earnings or news")
	}
	dividends := fixtures["finance/DividendTable"].Component.Props
	if len(dividends["Rows"].([]map[string]any)) != 4 || len(dividends["StatementRows"].([]map[string]any)) != 7 {
		t.Fatal("DividendTable lost source dividends or statement lines")
	}
}
