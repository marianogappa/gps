package main

import (
	"flag"
	"os"

	"github.com/marianogappa/gps/pkg"
)

func main() {
	separator := flag.String("separator", "tab", "Separator for output lines: tab or comma")
	flag.Parse()

	pkg.GPSFromReader(
		os.Stdin,
		pkg.WithOutputSeparator(*separator),
		pkg.WithOutputToWriter(os.Stdout),
	)
}
