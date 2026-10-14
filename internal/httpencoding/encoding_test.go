package httpencoding

import (
	"bytes"
	"compress/gzip"
	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"io"
	"net/http"
	"testing"
)

func TestNegotiate(t *testing.T) {
	for _, tc := range []struct {
		header, want string
		ok           bool
	}{
		{"", "identity", true}, {"gzip", "gzip", true}, {"gzip, br, zstd", "zstd", true},
		{"zstd;q=0.2, br;q=0.9, gzip;q=0.4", "br", true},
		{"gzip;q=0, br;q=0, zstd;q=0", "identity", true},
		{"identity;q=0, *;q=0", "", false}, {"*;q=0", "", false},
		{"*;q=0, gzip;q=1", "gzip", true}, {"*;q=1, zstd;q=0", "br", true},
		{"GZIP;Q=1", "gzip", true}, {"gzip;q=0.5, identity;q=1", "identity", true},
		{"gzip;q=0, gzip;q=1", "identity", true}, {"gzip;q=NaN", "identity", true},
		{"gzip;q=1.1", "identity", true}, {"gzip;q=.8", "identity", true},
		{"gzip;q=0.1234", "identity", true}, {"gzip;q=0.8;foo=bar", "identity", true},
		{"deflate", "identity", true}, {"deflate,identity;q=0", "", false},
		{"gzip;q=0.001", "gzip", true}, {"gzip;q=1.000", "gzip", true},
	} {
		t.Run(tc.header, func(t *testing.T) {
			got, ok := Negotiate(tc.header)
			if got != tc.want || ok != tc.ok {
				t.Fatalf("got %q %v; want %q %v", got, ok, tc.want, tc.ok)
			}
		})
	}
	if _, ok := NegotiateAvailable("gzip, identity;q=0", "identity"); ok {
		t.Fatal("identity exclusion ignored")
	}
}

func TestEncodingRoundTrip(t *testing.T) {
	data := bytes.Repeat([]byte("public text: hello world\n"), 100)
	for _, coding := range []string{"gzip", "br", "zstd", "identity"} {
		t.Run(coding, func(t *testing.T) {
			encoded, err := Encode(data, coding)
			if err != nil {
				t.Fatal(err)
			}
			var reader io.Reader = bytes.NewReader(encoded)
			switch coding {
			case "gzip":
				r, err := gzip.NewReader(reader)
				if err != nil {
					t.Fatal(err)
				}
				defer r.Close()
				reader = r
			case "br":
				reader = brotli.NewReader(reader)
			case "zstd":
				r, err := zstd.NewReader(reader)
				if err != nil {
					t.Fatal(err)
				}
				defer r.Close()
				reader = r
			}
			decoded, err := io.ReadAll(reader)
			if err != nil || !bytes.Equal(decoded, data) {
				t.Fatalf("roundtrip: %v", err)
			}
		})
	}
	if _, err := Encode(data, "deflate"); err == nil {
		t.Fatal("unsupported coding accepted")
	}
}

func TestAddVary(t *testing.T) {
	h := http.Header{"Vary": []string{"Origin, Accept-Language", "accept-encoding"}}
	AddVary(h, "Accept-Encoding")
	if len(h.Values("Vary")) != 2 {
		t.Fatal(h)
	}
	AddVary(h, "Cookie")
	if len(h.Values("Vary")) != 3 {
		t.Fatal(h)
	}
	h.Set("Vary", "*")
	AddVary(h, "Accept-Encoding")
	if h.Get("Vary") != "*" || len(h.Values("Vary")) != 1 {
		t.Fatal(h)
	}
}

func TestNormalizeOptions(t *testing.T) {
	for _, profile := range []string{"static", "dynamic", "request"} {
		o, err := Normalize(nil, profile)
		if err != nil || o.Disabled != (profile == "dynamic") {
			t.Fatalf("%s: %#v %v", profile, o, err)
		}
	}
	for _, o := range []Options{
		{Codecs: []string{"gzip", "gzip"}}, {Codecs: []string{"identity"}}, {Codecs: []string{"deflate"}},
		{Strategy: "random"}, {MinBytes: 64<<20 + 1}, {MinSavingsBytes: 64<<20 + 1}, {MinBytes: -1}, {MinSavingsBytes: -1}, {MinSavingsPercent: 101},
		{GzipLevel: 10}, {GzipLevel: -3}, {BrotliLevel: -1}, {BrotliLevel: 12}, {ZstdLevel: -6}, {ZstdLevel: 23},
	} {
		if _, err := Normalize(&o, "static"); err == nil {
			t.Fatalf("accepted invalid options %#v", o)
		}
	}
	if _, err := Normalize(nil, "unknown"); err == nil {
		t.Fatal("accepted unknown profile")
	}
	input := Options{Codecs: []string{"gzip"}}
	o, err := Normalize(&input, "dynamic")
	if err != nil {
		t.Fatal(err)
	}
	input.Codecs[0] = "br"
	if o.Codecs[0] != "gzip" || o.MinBytes != 0 || o.Disabled {
		t.Fatal("input alias or zero semantics", o)
	}
	empty, err := Normalize(&Options{Codecs: []string{}}, "static")
	if err != nil || len(empty.Codecs) != 0 {
		t.Fatal("empty list lost", empty, err)
	}
}

func TestSelectOptions(t *testing.T) {
	variants := map[string][]byte{"identity": make([]byte, 1000), "gzip": make([]byte, 100), "br": make([]byte, 50)}
	for _, tc := range []struct {
		strategy, header, want string
		threshold              int
		disabled               bool
		ok                     bool
	}{
		{"preference", "gzip,br", "gzip", 0, false, true},
		{"smallest", "gzip,br", "br", 0, false, true},
		{"smallest", "gzip;q=1,br;q=0.5", "gzip", 0, false, true},
		{"smallest", "gzip,br", "identity", 2000, false, true},
		{"smallest", "gzip,br,identity;q=0", "br", 2000, false, true},
		{"smallest", "*;q=0", "", 0, false, false},
		{"preference", "gzip,identity;q=0", "", 0, true, false},
		{"preference", "gzip", "identity", 0, true, true},
	} {
		o, _ := Normalize(&Options{Codecs: []string{"gzip", "br"}, Strategy: tc.strategy, MinBytes: tc.threshold, Disabled: tc.disabled}, "static")
		got, ok := Select(tc.header, variants, o)
		if got != tc.want || ok != tc.ok {
			t.Fatalf("%#v got %s %v", tc, got, ok)
		}
	}
	for _, o := range []Options{{MinBytes: 1001}, {MinSavingsBytes: 901}, {MinSavingsPercent: 91}} {
		if MeetsThresholds(1000, 100, o) {
			t.Fatal("threshold ignored", o)
		}
	}
	if !MeetsThresholds(1000, 1100, Options{}) {
		t.Fatal("zero thresholds must disable checks")
	}
}

func TestSelectEncodedOptions(t *testing.T) {
	data := bytes.Repeat([]byte("public dynamic response "), 100)
	for _, strategy := range []string{"preference", "smallest"} {
		for _, coding := range []string{"gzip", "br", "zstd"} {
			o, _ := Normalize(&Options{Codecs: []string{coding}, Strategy: strategy, GzipLevel: 1, BrotliLevel: 1, ZstdLevel: 1}, "dynamic")
			body, got, ok, err := SelectEncoded(data, coding, o)
			if err != nil || !ok || got != coding || len(body) >= len(data) {
				t.Fatalf("encode %s/%s: %s %v %v", strategy, coding, got, ok, err)
			}
			o.MinSavingsBytes = len(data)
			_, got, ok, err = SelectEncoded(data, coding, o)
			if err != nil || !ok || got != "identity" {
				t.Fatal("threshold fallback", got, ok, err)
			}
			_, got, ok, err = SelectEncoded(data, coding+",identity;q=0", o)
			if err != nil || !ok || got != coding {
				t.Fatal("forbidden identity fallback", got, ok, err)
			}
		}
	}
	// Invalid unselected encoder proves preference does not encode every codec.
	o := Options{Codecs: []string{"br", "gzip"}, Strategy: "preference", BrotliLevel: 1, GzipLevel: 99}
	if _, got, ok, err := SelectEncoded(data, "br,gzip", o); err != nil || !ok || got != "br" {
		t.Fatal(got, ok, err)
	}
}

func TestEncodingCustomLevelsRoundTrip(t *testing.T) {
	data := bytes.Repeat([]byte("public data with repeated fields and varying separators;\n"), 100)
	for _, tc := range []struct {
		coding string
		level  int
	}{{"gzip", -2}, {"gzip", -1}, {"gzip", 1}, {"gzip", 9}, {"br", 1}, {"br", 8}, {"zstd", -5}, {"zstd", 1}, {"zstd", 5}} {
		o, _ := Normalize(nil, "static")
		switch tc.coding {
		case "gzip":
			o.GzipLevel = tc.level
		case "br":
			o.BrotliLevel = tc.level
		case "zstd":
			o.ZstdLevel = tc.level
		}
		body, err := EncodeWithOptions(data, tc.coding, o)
		if err != nil {
			t.Fatal(err)
		}
		var r io.Reader = bytes.NewReader(body)
		switch tc.coding {
		case "gzip":
			decoder, e := gzip.NewReader(r)
			if e != nil {
				t.Fatal(e)
			}
			defer decoder.Close()
			r = decoder
		case "br":
			r = brotli.NewReader(r)
		case "zstd":
			decoder, e := zstd.NewReader(r)
			if e != nil {
				t.Fatal(e)
			}
			defer decoder.Close()
			r = decoder
		}
		decoded, err := io.ReadAll(r)
		if err != nil || !bytes.Equal(data, decoded) {
			t.Fatalf("%s level %d roundtrip failed: %v", tc.coding, tc.level, err)
		}
	}
}
