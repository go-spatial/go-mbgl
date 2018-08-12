package mbgl_test

import (
	"testing"
	"github.com/go-spatial/go-mbgl/mbgl"
	"github.com/arolek/p"
)

func TestDefaultFileSource(t *testing.T) {

	type tcase struct {
		cachePath, assetRoot string
		maxCache *uint64
	}

	fn := func(tc tcase, t *testing.T) {
		mbgl.NewDefaultFileSource(tc.cachePath, tc.assetRoot, tc.maxCache)

	}

	testcases := map[string]tcase{
		"1" : {
			cachePath: "",
			assetRoot: "https://osm.tegola.io/capabilities/osm",
			maxCache: nil,
		},
		"2" : {
			cachePath: "./cache",
			assetRoot: "https://osm.tegola.io/capabilities/osm",
			maxCache: p.Uint64(200),
		},
	}

	for k, v := range testcases {
		t.Run(k, func (t *testing.T) {
			fn(v, t)
		})
	}

}

func TestFileSourceInterface(t *testing.T) {
	fs := mbgl.NewDefaultFileSource("", "", nil)
	mbgl.FileSource(fs).Destruct()
}
