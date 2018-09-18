package mbgl

import "testing"

func TestMap(t *testing.T) {
	type tcase struct {
		size       Size
		pixelRatio float32
		src        FileSource
		sched      Scheduler
	}

	fn := func(tc tcase, t *testing.T) {
		front := NewHeadlessFrontend(tc.size,
			tc.pixelRatio,
			tc.src,
			tc.sched,
			nil, nil)

		NewMap(tc.front,
			tc.size,
			tc.pixelRatio,
			tc.src,
			tc.sched).
			Destruct()

		front.Destruct()
		tc.size.Destruct()
		tc.src.Destruct()
		tc.sched.Destruct()
	}

	testcases := map[string]tcase{
		"1": {
			size:       Size{256, 256},
			pixelRatio: 1.0,
			src:        NewDefaultFileSource("", "https://osm.tegola.io/maps/osm/style.json", 0),
			sched:      NewThreadPool(4),
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}
