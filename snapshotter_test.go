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
		ratio float32
	}

	fn := func(tc tcase, t *testing.T) {
		img := Snapshot(tc.src, tc.ext, tc.size, tc.ratio)

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

	tile := slippy.NewTileLatLon(12, 33, 117)

	exts["1"] = tile.Extent4326()

	z, x, y := tile.ZXY()
	tile = slippy.NewTile(z, x + 1, y)
	exts["2"] = tile.Extent4326()

	fmt.Println(exts["1"])

	testcases := map[string]tcase {
		"1": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(1000, 1000),
			ratio: 1.0,
		},
		"2": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(1000, 1000),
			ratio: 2.0,
		},
		"3": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(1000, 1000),
			ratio: 0.5,
		},
		"4": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(500, 500),
			ratio: 2.0,
		},
		"5": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(500, 500),
			ratio: 0.5,
		},
		"6": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(500, 500),
			ratio: 1.0,
		},
		"7": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(2000, 2000),
			ratio: 1.0,
		},
		"small": {
			src: "https://osm.tegola.io/maps/osm/style.json",
			ext: exts["1"],
			size: image.Pt(1, 1),
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}
