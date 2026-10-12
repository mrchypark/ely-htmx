package ui

import (
	"strings"
	"testing"
)

func TestCodesEncodeCallerPayloadWithoutHTMLInjection(t *testing.T) {
	first, err := qrSVG("연구 <source> & citation")
	if err != nil || !strings.Contains(string(first), "data-ely-code=") || !strings.Contains(string(first), "&lt;source&gt;") {
		t.Fatalf("QR encoding: %v %s", err, first)
	}
	second, err := qrSVG("Different caller payload")
	if err != nil || first == second {
		t.Fatalf("QR payload did not change the encoded geometry: %v", err)
	}
	if _, err := barcodeSVG("ELY HTMX-123"); err != nil {
		t.Fatal(err)
	}
	if _, err := barcodeSVG(strings.Repeat("x", 129)); err == nil {
		t.Fatal("oversized linear barcode accepted")
	}
	if _, err := qrSVG(strings.Repeat("x", 2049)); err == nil {
		t.Fatal("oversized QR payload accepted")
	}
}
