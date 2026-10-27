// Blog builds the publication around Margo and validates the generated site.
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: blog build|check|serve")
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "build":
		err = buildSite(".", time.Now())
	case "check":
		err = checkSite(".", "public", time.Now())
	case "serve":
		fmt.Println("Preview: http://127.0.0.1:1313")
		server := &http.Server{Addr: "127.0.0.1:1313", Handler: http.FileServer(http.Dir("public")), ReadHeaderTimeout: 5 * time.Second}
		err = server.ListenAndServe()
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
