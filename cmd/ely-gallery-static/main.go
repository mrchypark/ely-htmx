// Command ely-gallery-static exports the complete gallery for static hosting.
package main

import (
	"flag"
	"fmt"
	"log"

	gallery "github.com/mrchypark/ely-htmx/internal/gallery"
)

func main() {
	dir := flag.String("dir", "dist", "static output directory")
	base := flag.String("base", "/ely-htmx/", "absolute deployment URL path")
	flag.Parse()
	result, err := gallery.ExportStatic(*dir, *base)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exported %d stories and %d pages to %s (base %s)\n", result.Stories, len(result.Pages), *dir, result.Base)
}
