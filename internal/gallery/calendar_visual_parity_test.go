package uiexamples

import "testing"

func TestCalendarTimeFixtureHasWholeDayAndIsolatedOverlay(t *testing.T) {
	week, day := calendarWeekTimeGrid(), calendarDayTimeGrid()
	if len(week) != 48 || len(day) != 48 {
		t.Fatalf("half-hour rows: week %d day %d", len(week), len(day))
	}
	if week[0]["Time"] != "00:00" || week[47]["Time"] != "23:30" {
		t.Fatal("time grid omits part of the day")
	}
	for i, row := range week {
		if _, ok := row["OverlayTime"]; ok {
			t.Fatal("week duplicates clock labels in its single gutter")
		}
		if day[i]["OverlayTime"] != day[i]["Time"] {
			t.Fatal("day overlay lost its clock label")
		}
		if len(day[i]["Cells"].([]map[string]any)) != 1 {
			t.Fatal("day has multiple columns")
		}
	}
	events := week[42]["Cells"].([]map[string]any)[6]["Events"].([]map[string]any)
	if events[0]["Value"] != "night" || events[0]["Minutes"] != 180 || events[0]["End"] != "2026-10-19T02:00" {
		t.Fatal("cross-midnight caller duration or original end lost")
	}
}
