package mbgl

import (
	"testing"
	"github.com/go-spatial/geom"
	"image"
	"os"
	"fmt"
	"path/filepath"
	"image/png"
	"github.com/go-spatial/geom/slippy"
	"strings"
)

func TestSnapshot(t *testing.T) {
	type tcase struct {
		src string
		ext *geom.Extent
		size image.Point
	}

	fn := func(tc tcase, t *testing.T) {
		img := Snapshot(tc.src, tc.ext, tc.size)

		fname := os.DevNull
		if evar := os.Getenv("MBGL_TEST_OUT_DIR"); evar != "" {
			fmt.Println("outputing to ",evar)
			os.MkdirAll(evar, 0600)

			fname = strings.Replace(t.Name(), "/", "-", -1)
			fname = filepath.Join(evar,  fname + ".png")
		}

		f, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}

	exts := make(map[string]*geom.Extent)

	exts["1"] = slippy.NewTileLatLon(12, 33, 117, 0, geom.WebMercator).Extent(geom.WGS84)

	fmt.Println(exts["1"])

	testcases := map[string]tcase {
		"1": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(1000, 1000),
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}
