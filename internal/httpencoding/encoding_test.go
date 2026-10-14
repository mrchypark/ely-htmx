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
