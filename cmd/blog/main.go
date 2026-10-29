// Blog is the Guigoes entry point for the reusable publication engine.
package main

import (
	"fmt"
	blog "github.com/araihu/gazeta"
	"os"
)

func main() {
	if err := blog.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
