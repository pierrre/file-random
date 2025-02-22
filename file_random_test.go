package filerandom

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/pierrre/assert"
)

func Test(t *testing.T) {
	ctx := t.Context()
	fsys := fstest.MapFS{
		"empty": &fstest.MapFile{},
		"small": &fstest.MapFile{
			Data: []byte("a"),
		},
		"large": &fstest.MapFile{
			Data: []byte("aaaaa"),
		},
	}
	fps, err := Get(ctx, WithFSs([]fs.FS{fsys}), WithMinSize(2))
	assert.NoError(t, err)
	assert.SliceLen(t, fps, 1)
	fp := fps.GetRandom()
	assert.Equal(t, fp.Path, "large")
}

func TestFilesGetRandomEmpty(t *testing.T) {
	var fps Files
	fp := fps.GetRandom()
	assert.Zero(t, fp)
}
