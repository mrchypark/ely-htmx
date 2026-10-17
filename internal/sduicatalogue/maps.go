package sduicatalogue

import (
	"fmt"
	"strings"
)

// Map coordinates are real numbers even when a fixture happens to contain an
// integer (or omits an optional coordinate). Limit this reviewed annotation to
// the six coordinate sinks of map-view, including its qualified aliases. Other
// numeric props, coordinate display strings and GeoJSON keep their own schemas.
func applyMapCoordinateTypes(c *Contract) error {
	if c.Template != "map-view" {
		if c.Chapter != "maps" {
			return nil
		}
		switch c.Primary {
		case "MapView", "LocationPicker", "RouteLine", "GeoHeatmap", "MapCluster", "ChoroplethMap":
		default:
			return nil
		}
	}
	for _, path := range []string{
		"center.latitude", "center.longitude",
		"points.*.latitude", "points.*.longitude",
		"routes.*.points.*.latitude", "routes.*.points.*.longitude",
	} {
		root := Property{Type: "object", Properties: c.Properties}
		if err := mapCoordinateType(&root, strings.Split(path, ".")); err != nil {
			return fmt.Errorf("map coordinate %s: %w", path, err)
		}
	}
	return nil
}

func mapCoordinateType(p *Property, path []string) error {
	if len(path) == 0 {
		switch p.Type {
		case "number", "integer", "string":
			p.Type, p.Format = "number", ""
			return nil
		}
		return fmt.Errorf("unexpected coordinate type %q", p.Type)
	}
	if path[0] == "*" && p.Type == "array" && p.Items != nil {
		return mapCoordinateType(p.Items, path[1:])
	}
	child, ok := p.Properties[path[0]]
	if !ok || p.Type != "object" {
		return fmt.Errorf("missing coordinate field %q", path[0])
	}
	if err := mapCoordinateType(&child, path[1:]); err != nil {
		return err
	}
	p.Properties[path[0]] = child
	return nil
}
