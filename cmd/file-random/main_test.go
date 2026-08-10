package main

import (
	"bytes"
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	filerandom "github.com/pierrre/file-random"
)

func TestOK(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err = run(ctx, fl, stdout, l, nil, nil)
	assert.NoError(t, err)
	expectedStdout := filepath.Join(wd, "testdata", "large") + "\n"
	assert.Equal(t, stdout.String(), expectedStdout)
	assert.Zero(t, stderr.String())
}

func TestOpenFile(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	openFileCalled := false
	expectedPath := filepath.Join(wd, "testdata", "large")
	openFile := func(p string) error {
		openFileCalled = true
		assert.Equal(t, p, expectedPath)
		return nil
	}
	err = run(ctx, fl, stdout, l, openFile, nil)
	assert.NoError(t, err)
	expectedStdout := expectedPath + "\n"
	assert.Equal(t, stdout.String(), expectedStdout)
	assert.Zero(t, stderr.String())
	assert.True(t, openFileCalled)
}

func TestLoop(t *testing.T) {
	ctx := t.Context()
	ctx, cancel := context.WithCancel(ctx)
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.loop = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	waitEnter := func() {
		cancel()
	}
	err = run(ctx, fl, stdout, l, nil, waitEnter)
	assert.NoError(t, err)
	expectedStdout := filepath.Join(wd, "testdata", "large") + "\n"
	assert.Equal(t, stdout.String(), expectedStdout)
	assert.Zero(t, stderr.String())
}

func TestErrorOpenFile(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	openFile := func(p string) error {
		return errors.New("error")
	}
	err = run(ctx, fl, stdout, l, openFile, nil)
	assert.Error(t, err)
	assert.Zero(t, stderr.String())
}

func TestErrorOpenFileContinue(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.continueOnError = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	openFile := func(p string) error {
		return errors.New("error")
	}
	err = run(ctx, fl, stdout, l, openFile, nil)
	assert.NoError(t, err)
	assert.NotZero(t, stderr.String())
}

func TestErrorNoRoots(t *testing.T) {
	ctx := t.Context()
	fl := newFlags()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	err := run(ctx, fl, stdout, l, nil, nil)
	assert.Error(t, err)
	assert.Zero(t, stdout.String())
	assert.Zero(t, stderr.String())
}

func TestBuildFSsRootSlash(t *testing.T) {
	fl := newFlags()
	fl.roots = []string{"/"}
	fsyss := buildFSs(fl)
	assert.SliceLen(t, fsyss, 1)
	entries, err := fs.ReadDir(fsyss[0], ".")
	assert.NoError(t, err)
	assert.NotZero(t, len(entries))
}

func TestBuildOptionsMinSizeZero(t *testing.T) {
	ctx := t.Context()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 0
	fl.roots = []string{path.Join(wd, "testdata")}
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr, nil))
	optfs := buildOptions(fl, l)
	fps, err := filerandom.Get(ctx, optfs...)
	assert.NoError(t, err)
	assert.SliceLen(t, fps, 3)
	assert.Zero(t, stderr.String())
}
