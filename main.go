package main

import (
	"fmt"
	"os"

	"buf.build/go/bufplugin/check"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Println(version)
			os.Exit(0)
		}
	}
	check.Main(buildSpec())
}
