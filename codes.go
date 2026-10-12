package ui

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
)

func qrSVG(value any) (template.HTML, error) {
	return encodedSVG(value, true)
}

func barcodeSVG(value any) (template.HTML, error) {
	return encodedSVG(value, false)
}

func encodedSVG(value any, matrix bool) (template.HTML, error) {
	if value == nil || fmt.Sprint(value) == "" {
		return `<span data-ely-empty>Provide a value to encode.</span>`, nil
	}
	content := fmt.Sprint(value)
	if len(content) > 2048 || (!matrix && len(content) > 128) {
		return "", fmt.Errorf("barcode content exceeds the supported length")
	}
	var code barcode.Barcode
	var err error
	quiet, height, label := 10, 48, "Code 128 barcode"
	if matrix {
		code, err = qr.Encode(content, qr.M, qr.Auto)
		quiet, label = 4, "QR code"
	} else {
		code, err = code128.Encode(content)
	}
	if err != nil {
		return "", fmt.Errorf("encode %s: %w", label, err)
	}
	bounds := code.Bounds()
	if matrix {
		height = bounds.Dy()
	}
	var path strings.Builder
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			r, g, b, _ := code.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			if r+g+b < 3*32768 {
				h := 1
				if !matrix {
					h = height
				}
				fmt.Fprintf(&path, "M%d %dh1v%dh-1z", x+quiet, y+quiet, h)
			}
		}
	}
	return template.HTML(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" role="img" aria-label="%s" data-ely-code="%s" shape-rendering="crispEdges"><title>%s</title><rect width="100%%" height="100%%" fill="white"/><path d="%s" fill="black"/></svg>`, bounds.Dx()+quiet*2, height+quiet*2, label, template.HTMLEscapeString(content), template.HTMLEscapeString(content), path.String())), nil
}
