package mbgl

import (
	"testing"
	"image/png"
	"os"
	"fmt"
)

func TestNewMapSnapshotter(t *testing.T) {
	type tcase struct {
		src        FileSource
		sched      Scheduler
		style      string
		size       Size
		pixelRatio float32
		camOpts    *CameraOptions
		region     *LatLngBounds
		cacheDir   *string
	}

	fn := func(tc tcase, t *testing.T) {
		ms := NewMapSnapshotter(
			tc.src,
			tc.sched,
			tc.style,
			tc.size,
			tc.pixelRatio,
			tc.camOpts,
			tc.region,
			tc.cacheDir)

		ms.Destruct()
	}

	testcases := map[string]tcase{
		"1": {
			src:        NewDefaultFileSource("", "", nil),
			sched:      NewThreadPool(4),
			style:      "https://osm.tegola.io/maps/osm/style.json",
			size:       Size{Width: 100, Height: 100},
			pixelRatio: 1.0,
			camOpts:    nil,
			region:     nil,
			cacheDir:   nil,
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}

func TestSnapshot(t *testing.T) {
	type tcase struct {
		ms *MapSnapshotter
	}

	fn := func(tc tcase, t *testing.T) {
		cImg := tc.ms.Snapshot()
		img := cImg.Image()

		fname := os.DevNull
		if evar := os.Getenv("MBGL_TEST_ON_DISK"); evar != "" {
			fmt.Println("outputing to ",evar)
			fname = evar
		}

		f, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE, 0600)
		if err != nil {
			t.Fatalf("unexpected errro %v", err)
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			t.Fatalf("unexpected errro %v", err)
		}

		tc.ms.Destruct()
	}

	tpool := NewThreadPool(4)
	SchedulerSetCurrent(tpool)

	testcases := map[string]tcase{
		"1": {
			ms: NewMapSnapshotter(
				NewDefaultFileSource("", "", nil),
				tpool,
				"https://osm.tegola.io/maps/osm/style.json",
				Size{Height: 100, Width: 100},
				1.0,
				nil,
				nil,
				nil),
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}
