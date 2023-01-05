package main

import (
	"flag"
	"os"
)

type flags struct {
	verbose         bool
	continueOnError bool
	minSize         int64
	open            bool
	loop            bool
	roots           []string
}

func parseFlags() *flags {
	fl := newFlags()
	fls := flag.NewFlagSet("", flag.ExitOnError)
	fls.BoolVar(&fl.verbose, "v", fl.verbose, "verbose")
	fls.BoolVar(&fl.continueOnError, "continue-on-error", fl.continueOnError, "continue on error")
	fls.Int64Var(&fl.minSize, "min-size", fl.minSize, "min size")
	fls.BoolVar(&fl.open, "open", fl.open, "open")
	fls.BoolVar(&fl.loop, "loop", fl.loop, "loop")
	_ = fls.Parse(os.Args[1:])
	fl.roots = fls.Args()
	return fl
}

func newFlags() *flags {
	return &flags{
		minSize: 1,
	}
}
