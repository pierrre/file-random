package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/pierrre/errors"
)

func TestOK(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fl := newFlags()
	fl.minSize = 2
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	err = run(ctx, fl, stdout, l, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	expectedStdout := filepath.Join(wd, "testdata", "large") + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("unexpected stdout: got %q, want %q", stdout.String(), expectedStdout)
	}
	if stderr.String() != "" {
		t.Fatalf("unexpected stderr: got %q, want %q", stderr.String(), "")
	}
}

func TestOpenFile(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	openFileCalled := false
	expectedPath := filepath.Join(wd, "testdata", "large")
	openFile := func(p string) error {
		openFileCalled = true
		if p != expectedPath {
			t.Fatalf("unexpected path: got %q, want %q", p, expectedPath)
		}
		return nil
	}
	err = run(ctx, fl, stdout, l, openFile, nil)
	if err != nil {
		t.Fatal(err)
	}
	expectedStdout := expectedPath + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("unexpected stdout: got %q, want %q", stdout.String(), expectedStdout)
	}
	if stderr.String() != "" {
		t.Fatalf("unexpected stderr: got %q, want %q", stderr.String(), "")
	}
	if !openFileCalled {
		t.Fatal("openFile not called")
	}
}

func TestLoop(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fl := newFlags()
	fl.minSize = 2
	fl.loop = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	waitEnter := func() {
		cancel()
	}
	err = run(ctx, fl, stdout, l, nil, waitEnter)
	if err != nil {
		t.Fatal(err)
	}
	expectedStdout := filepath.Join(wd, "testdata", "large") + "\n"
	if stdout.String() != expectedStdout {
		t.Fatalf("unexpected stdout: got %q, want %q", stdout.String(), expectedStdout)
	}
	if stderr.String() != "" {
		t.Fatalf("unexpected stderr: got %q, want %q", stderr.String(), "")
	}
}

func TestErrorOpenFile(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	openFile := func(p string) error {
		return errors.New("error")
	}
	err = run(ctx, fl, stdout, l, openFile, nil)
	if err == nil {
		t.Fatal("no error")
	}
	if stderr.String() != "" {
		t.Fatalf("unexpected stderr: got %q, want %q", stderr.String(), "")
	}
}

func TestErrorOpenFileContinue(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fl := newFlags()
	fl.minSize = 2
	fl.open = true
	fl.continueOnError = true
	fl.roots = []string{path.Join(wd, "testdata")}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	l := log.New(stderr, "", 0)
	openFile := func(p string) error {
		return errors.New("error")
	}
	err = run(ctx, fl, stdout, l, openFile, nil)
	if err != nil {
		t.Fatal(err)
	}
	if stderr.Len() == 0 {
		t.Fatal("no error log")
	}
}
