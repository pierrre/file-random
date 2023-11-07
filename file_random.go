// Package filerandom provides utilities to get random files.
package filerandom

import (
	"context"
	"io/fs"
	"math/rand"

	"github.com/pierrre/errors"
)

type options struct {
	fss          []fs.FS
	minSize      int64
	errorHandler func(context.Context, error)
}

func newOptions(optfs ...Option) *options {
	opts := &options{
		minSize: 1,
	}
	for _, optf := range optfs {
		optf(opts)
	}
	return opts
}

// Option represents an option.
type Option func(*options)

// WithFSs is an [Option] that defines the filesystems to scan.
func WithFSs(fsyss []fs.FS) Option {
	return func(o *options) {
		o.fss = fsyss
	}
}

// WithMinSize is an [Option] that defines the minimum file size to consider.
func WithMinSize(minSize int64) Option {
	return func(o *options) {
		o.minSize = minSize
	}
}

// WithErrorHandler is an [Option] that defines the error handler.
//
// If it is defined, the error handler is called for each error, otherwise the error is returned.
func WithErrorHandler(f func(context.Context, error)) Option {
	return func(o *options) {
		o.errorHandler = f
	}
}

// Files represents a list of files.
type Files []*File

// GetRandom returns a random file.
//
// It returns nil if the list of files is empty.
func (fps Files) GetRandom() *File {
	if len(fps) == 0 {
		return nil
	}
	i := rand.Intn(len(fps)) //nolint:gosec // It is ok to use a non-crypto random here.
	return fps[i]
}

// File represents a file.
type File struct {
	// FSIndex is the index of the filesystem where the file is located.
	FSIndex int
	// Path is the path of the file in the filesystem.
	Path string
}

// Get returns a [Files].
func Get(ctx context.Context, optfs ...Option) (Files, error) {
	opts := newOptions(optfs...)
	fps, err := getFiles(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "get files")
	}
	return fps, nil
}

func getFiles(ctx context.Context, opts *options) ([]*File, error) {
	var res Files
	for fsysIdx, fsys := range opts.fss {
		wdf := newWalkDirFunc(ctx, opts, &res, fsysIdx)
		err := fs.WalkDir(fsys, ".", wdf)
		if err != nil {
			return nil, errors.Wrap(err, "walk dir")
		}
	}
	return res, nil
}

func newWalkDirFunc(ctx context.Context, opts *options, res *Files, fsysIdx int) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if opts.errorHandler != nil {
				err = errors.Wrap(err, "walk dir")
				opts.errorHandler(ctx, err)
				return nil
			}
			return errors.Wrap(err, "")
		}
		if !d.Type().IsRegular() {
			return nil
		}
		fi, err := d.Info()
		if err != nil {
			err = errors.Wrap(err, "info")
			if opts.errorHandler != nil {
				opts.errorHandler(ctx, err)
				return nil
			}
			return err
		}
		size := fi.Size()
		if size < opts.minSize {
			return nil
		}
		fp := &File{
			FSIndex: fsysIdx,
			Path:    path,
		}
		*res = append(*res, fp)
		return nil
	}
}
