package main

import (
	"bytes"
	"context"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/ext/davecghspew"
	"github.com/pierrre/assert/ext/pierrrecompare"
	"github.com/pierrre/assert/ext/pierrreerrors"
	"github.com/pierrre/errors"
	"golang.org/x/exp/slog"
)

func init() {
	pierrrecompare.Configure()
	davecghspew.ConfigureDefault()
	pierrreerrors.Configure()
}

func TestOK(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr))
	err = run(ctx, fl, stdout, l, nil, nil)
	assert.NoError(t, err)
	expectedStdout := filepath.Join(wd, "testdata", "large") + "\n"
	assert.Equal(t, stdout.String(), expectedStdout)
	assert.StringEmpty(t, stderr.String())
}

func TestOpenFile(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr))
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
	assert.StringEmpty(t, stderr.String())
	assert.True(t, openFileCalled)
}

func TestLoop(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.loop = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr))
	waitEnter := func() {
		cancel()
	}
	err = run(ctx, fl, stdout, l, nil, waitEnter)
	assert.NoError(t, err)
	expectedStdout := filepath.Join(wd, "testdata", "large") + "\n"
	assert.Equal(t, stdout.String(), expectedStdout)
	assert.StringEmpty(t, stderr.String())
}

func TestErrorOpenFile(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr))
	openFile := func(p string) error {
		return errors.New("error")
	}
	err = run(ctx, fl, stdout, l, openFile, nil)
	assert.Error(t, err)
	assert.StringEmpty(t, stderr.String())
}

func TestErrorOpenFileContinue(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.continueOnError = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(stderr))
	openFile := func(p string) error {
		return errors.New("error")
	}
	err = run(ctx, fl, stdout, l, openFile, nil)
	assert.NoError(t, err)
	assert.StringNotEmpty(t, stderr.String())
}
