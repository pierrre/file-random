// Package file-duplicate provides a command line tool to find duplicate files.
package main

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errverbose"
	filerandom "github.com/pierrre/file-random"
	"github.com/pkg/browser"
)

func main() {
	ctx := context.Background()
	fl := parseFlags()
	l := log.Default()
	br := bufio.NewReader(os.Stdin)
	waitEnter := func() {
		l.Println("Press enter to continue")
		_, _ = br.ReadString('\n')
	}
	err := run(ctx, fl, os.Stdout, l, browser.OpenFile, waitEnter)
	if err != nil {
		handleError(l.Fatalf, err)
	}
}

func run(ctx context.Context, fl *flags, w io.Writer, l *log.Logger, openFile func(p string) error, waitEnter func()) error {
	optfs := buildOptions(fl, l)
	fps, err := filerandom.Get(optfs...)
	if err != nil {
		return errors.Wrap(err, "get")
	}
	if len(fps) == 0 {
		return errors.New("no file")
	}
	for ctx.Err() == nil {
		fp := fps.GetRandom()
		root := fl.roots[fp.FSIndex]
		p := filepath.Join(root, fp.Path)
		_, _ = io.WriteString(w, p)
		_, _ = io.WriteString(w, "\n")
		if fl.open {
			err = openFile(p)
			if err != nil {
				err = errors.Wrap(err, "open file")
				if !fl.continueOnError {
					return err
				}
				handleError(l.Printf, err)
			}
		}
		if !fl.loop {
			break
		}
		waitEnter()
	}
	return nil
}

func buildOptions(fl *flags, l *log.Logger) []filerandom.Option {
	var optfs []filerandom.Option
	fsyss := make([]fs.FS, len(fl.roots))
	for i, root := range fl.roots {
		root = filepath.Clean(root)
		if root == "/" {
			root = ""
		}
		fsyss[i] = os.DirFS(root)
	}
	optfs = append(optfs, filerandom.WithFSs(fsyss))
	if fl.minSize != 0 {
		optfs = append(optfs, filerandom.WithMinSize(fl.minSize))
	}
	if fl.continueOnError {
		optfs = append(optfs, filerandom.WithErrorHandler(func(err error) {
			if fl.verbose {
				handleError(l.Printf, err)
			}
		}))
	}
	return optfs
}

func handleError(lf func(format string, v ...any), err error) {
	lf("Error: %v", errverbose.Formatter(err))
}
