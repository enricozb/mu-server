package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/enricozb/mu-server/api"
	"github.com/enricozb/mu-server/library"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fatalf("usage: mu-server [dir]")
	}

	dir := args[0]

	lib := library.New(dir)
	if err := lib.Init(); err != nil {
		fatalf("init library: %v", err)
	}

	if err := api.New(lib).Run(); err != nil {
		fatalf("run api: %v", err)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}
