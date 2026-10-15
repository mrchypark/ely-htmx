package main

import (
	"bytes"
	"testing"
)

func TestFixturesDeterministic(t *testing.T) {
	a, b := datasets(), datasets()
	if len(a) != 10 {
		t.Fatal(len(a))
	}
	for i := range a {
		if !bytes.Equal(a[i].data, b[i].data) || len(a[i].data) == 0 {
			t.Fatal(a[i].name)
		}
	}
}
func TestCodecRoundTrips(t *testing.T) {
	for _, d := range datasets() {
		for _, p := range profiles {
			if p.name != "gzip1" && p.name != "br1" && p.name != "zstd1" {
				continue
			}
			t.Run(d.name+"/"+p.name, func(t *testing.T) {
				b, e := encode(d.data, p)
				if e != nil {
					t.Fatal(e)
				}
				got, e := decode(b, p.coding)
				if e != nil || !bytes.Equal(got, d.data) {
					t.Fatalf("round trip: %v", e)
				}
			})
		}
	}
}
