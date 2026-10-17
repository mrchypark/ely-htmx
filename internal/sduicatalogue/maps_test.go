package sduicatalogue

import "testing"

func TestMapCoordinatesOverrideFixtureInferenceOnly(t *testing.T) {
	coordinate := func(typ string) Property {
		return Property{Type: "object", Properties: map[string]Property{
			"latitude":  {Type: typ, Target: "Latitude"},
			"longitude": {Type: typ, Target: "Longitude"},
		}}
	}
	point := coordinate("integer")
	routePoint := coordinate("string")
	route := Property{Type: "object", Properties: map[string]Property{"points": {Type: "array", Items: &routePoint}}}
	for _, contract := range []Contract{
		{Name: "basic/map-view", Template: "map-view"},
		{Name: "maps/MapMarker", Chapter: "maps", Primary: "MapView"},
		{Name: "maps/GeoJSON Layer", Chapter: "maps", Primary: "RouteLine"},
	} {
		t.Run(contract.Name, func(t *testing.T) {
			contract.Properties = map[string]Property{
				"center": coordinate("string"), "points": {Type: "array", Items: &point},
				"routes": {Type: "array", Items: &route}, "zoom": {Type: "integer"},
			}
			if err := applyMapCoordinateTypes(&contract); err != nil {
				t.Fatal(err)
			}
			for _, record := range []Property{contract.Properties["center"], *contract.Properties["points"].Items, *contract.Properties["routes"].Items.Properties["points"].Items} {
				for _, key := range []string{"latitude", "longitude"} {
					p := record.Properties[key]
					if p.Type != "number" || p.Format != "" || p.Target == "" || p.Required || p.Nullable {
						t.Fatalf("coordinate metadata changed: %+v", p)
					}
				}
			}
			if contract.Properties["zoom"].Type != "integer" {
				t.Fatal("unrelated integer widened")
			}
		})
	}
	for _, c := range []Contract{{Name: "other/MapView", Chapter: "other", Primary: "MapView"}, {Name: "maps/CoordinateDisplay", Chapter: "maps", Primary: "CoordinateDisplay"}} {
		c.Properties = map[string]Property{"center": coordinate("integer")}
		if err := applyMapCoordinateTypes(&c); err != nil {
			t.Fatal(err)
		}
		if c.Properties["center"].Properties["latitude"].Type != "integer" {
			t.Fatalf("annotation escaped map-view: %s", c.Name)
		}
	}
	if err := applyMapCoordinateTypes(&Contract{Template: "map-view"}); err == nil {
		t.Fatal("missing renderer coordinate schema was silently accepted")
	}
	v, err := canonicalValue(map[string]any{"Rows": 3, "Latitude": 20.0})
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range infer(v).Properties {
		if p.Type != "integer" {
			t.Fatal("global integer inference changed")
		}
	}
}
