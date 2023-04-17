package filerandom

import (
	"context"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/ext/davecghspew"
	"github.com/pierrre/assert/ext/pierrrecompare"
	"github.com/pierrre/assert/ext/pierrreerrors"
)

func init() {
	pierrrecompare.Configure()
	davecghspew.ConfigureDefault()
	pierrreerrors.Configure()
}

func Test(t *testing.T) {
	ctx := context.Background()
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
