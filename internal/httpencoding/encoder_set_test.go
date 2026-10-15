package httpencoding

import (
	"bytes"
	"testing"
)

func TestEncoderSetMatchesFreshAndOwnsOutputs(t *testing.T) {
	for _, profile := range []string{"static", "dynamic"} {
		options, _ := Normalize(nil, profile)
		set := NewEncoderSet(options)
		var saved, expected [][]byte
		for _, data := range [][]byte{bytes.Repeat([]byte("first public asset."), 300), []byte("tiny"), bytes.Repeat([]byte("second asset with different data."), 1000), {}, bytes.Repeat([]byte("first public asset."), 300)} {
			for _, coding := range []string{"gzip", "br", "zstd", "identity"} {
				got, err := set.Encode(data, coding)
				if err != nil {
					t.Fatal(err)
				}
				want, err := EncodeWithOptions(data, coding, options)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, want) {
					t.Fatalf("%s/%s batch output differs from fresh encoder", profile, coding)
				}
				saved = append(saved, got)
				expected = append(expected, bytes.Clone(want))
			}
			// Also verify identity output does not borrow the caller's input buffer.
			for i := range data {
				data[i] = 'x'
			}
		}
		set.Close()
		set.Close()
		for i := range saved {
			if !bytes.Equal(saved[i], expected[i]) {
				t.Fatalf("%s result %d aliases mutable workspace", profile, i)
			}
		}
		if _, err := set.Encode([]byte("x"), "gzip"); err == nil {
			t.Fatal("closed encoder accepted input")
		}
	}
}

func TestEncoderSetInvalidCodec(t *testing.T) {
	options, _ := Normalize(nil, "static")
	set := NewEncoderSet(options)
	defer set.Close()
	if _, err := set.Encode([]byte("x"), "unknown"); err == nil {
		t.Fatal("accepted invalid codec")
	}
	if _, err := set.Encode([]byte("x"), "gzip"); err != nil {
		t.Fatal("invalid codec damaged batch", err)
	}
}
