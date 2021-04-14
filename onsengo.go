// List and browse onsen.ag radio shows.
package main

import (
	"os"

	"github.com/adios/onsengo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
