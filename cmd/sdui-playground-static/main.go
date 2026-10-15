// Command sdui-playground-static exports the serverless SDUI editor.
package main

import (
	"flag"
	"github.com/mrchypark/ely-htmx/internal/playground"
	"log"
)

func main() {
	dir := flag.String("dir", "dist/playground", "output directory")
	flag.Parse()
	if err := playground.ExportStatic(*dir); err != nil {
		log.Fatal(err)
	}
}
