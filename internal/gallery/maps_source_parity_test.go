package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMapsSourceStories(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]string{
		"mapview-mapmarker-mappopup": {`data-map-latitude="20"`, `data-map-longitude="10"`, `data-map-initial-zoom="1.2"`, `data-map-width="720"`, `data-map-height="360"`, `data-key="london"`, `data-lat="51.5072"`, `data-lon="-0.1276"`, `Land: Natural Earth`},
		"coordinatedisplay":          {"33.9249° S, 18.4241° E", "33°55′30″ S, 18°25′27″ E", "degrees, minutes, seconds", "decimal", "data-map-copy"},
		"locationpicker":             {`data-map-latitude="48.8566"`, `data-map-longitude="2.3522"`, `data-map-initial-zoom="3"`, `data-map-width="480"`, `data-map-height="288"`, `data-map-format="dms"`},
		"routeline-geojson-layer":    {`data-map-latitude="38"`, `data-map-longitude="-40"`, `data-map-initial-zoom="2.3"`, `The Bermuda Triangle`, `A ring around the Azores`, `data-lat="64.1466"`},
		"geoheatmap":                 {`data-map-latitude="20"`, `data-map-longitude="30"`, `data-map-initial-zoom="1"`, `data-map-hue="6"`},
		"mapcluster":                 {`data-map-latitude="42"`, `data-map-longitude="5"`, `data-map-initial-zoom="3"`, `data-key="pin-41"`},
		"choroplethmap-worldmap":     {`data-map-latitude="25"`, `data-map-initial-zoom="1.4"`, `data-map-initial-zoom="3.2"`, `Countries: Natural Earth`, `data-value="120400"`, `data-value="540"`},
	}
	for path, wants := range cases {
		t.Run(path, func(t *testing.T) {
			r := httptest.NewRecorder()
			s.ServeHTTP(r, httptest.NewRequest("GET", "/components/maps/"+path, nil))
			if r.Code != 200 {
				t.Fatalf("status %d", r.Code)
			}
			body := r.Body.String()
			for _, want := range wants {
				if !strings.Contains(body, want) {
					t.Errorf("missing %q", want)
				}
			}
			if strings.Contains(body, "37.5665,126.9780") {
				t.Error("legacy unrelated Seoul coordinate remains")
			}
			if path == "choroplethmap-worldmap" && strings.Index(body, `data-map-mode="world"`) > strings.Index(body, `data-map-mode="choropleth"`) {
				t.Error("WorldMap must precede rainfall")
			}
		})
	}
}

func TestMapsSeededSourceAndAccountAliasFixtures(t *testing.T) {
	cases := advancedFormsCases()
	found := map[string]bool{}
	for _, item := range cases {
		key := item.Component.Chapter + "/" + item.Component.Name
		p := item.Component.Props
		switch key {
		case "maps/GeoHeatmap":
			points := p["Points"].([]map[string]any)
			if len(points) != 150 {
				t.Fatalf("heat points %d", len(points))
			}
			next := mapSourceNoise(41)
			want := 35.68 + (next()-.5)*6
			if points[0]["Latitude"] != want {
				t.Error("heat point did not use source seed")
			}
		case "maps/MapCluster":
			if len(p["Points"].([]map[string]any)) != 42 {
				t.Error("cluster source count differs")
			}
		case "account/PaymentMethodForm", "account/InvitationList", "onboarding/ContactSupport":
			found[key] = true
			if len(p["Fields"].([]map[string]any)) == 0 {
				t.Errorf("%s alias lost fields", key)
			}
		}
	}
	for _, key := range []string{"account/PaymentMethodForm", "account/InvitationList", "onboarding/ContactSupport"} {
		if !found[key] {
			t.Errorf("missing dedicated alias %s", key)
		}
	}
}
