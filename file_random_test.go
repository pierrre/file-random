package filerandom

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func Test(t *testing.T) {
	fsys := fstest.MapFS{
		"empty": &fstest.MapFile{},
		"small": &fstest.MapFile{
			Data: []byte("a"),
		},
		"large": &fstest.MapFile{
			Data: []byte("aaaaa"),
		},
	}
	fps, err := Get(WithFSs([]fs.FS{fsys}), WithMinSize(2))
	if err != nil {
		t.Fatal(err)
	}
	if len(fps) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(fps), 1)
	}
	fp := fps.GetRandom()
	if fp.Path != "large" {
		t.Fatalf("unexpected path: got %q, want %q", fp.Path, "large")
	}
}

func TestFilesGetRandomEmpty(t *testing.T) {
	var fps Files
	fp := fps.GetRandom()
	if fp != nil {
		t.Fatalf("unexpected file: got %v, want nil", fp)
	}
}
