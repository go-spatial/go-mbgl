package mbgl

import (
	"testing"
	"github.com/arolek/p"
)

func TestCameraOptions(t *testing.T) {
	type coptions struct {
		center *LatLng
		padding *EdgeInsets
		anchor *Point
		zoom *float64
		angle *float64
		pitch *float64
	}

	type tcase struct {
		opts []coptions
	}

	fn := func (tc tcase, t *testing.T) {

		// this test will leak memory out of
		opts := &CameraOptions{}
		opts.update()
		for _, v:= range tc.opts {
			opts.Center = v.center
			opts.Padding = v.padding
			opts.Anchor = v.anchor
			opts.Zoom = v.zoom
			opts.Angle = v.angle
			opts.Pitch = v.pitch
			opts.update()
		}

		opts.Destruct()
	}

	testcases := map[string]tcase{
		"1": {
			[]coptions{
				{},
			},
		},

		"2": {
			[]coptions{
				{},
				{
					center: NewLatLng(32, 117),
					padding: NewEdgeInsets(0, 0, 0, 0),
					anchor: &Point{X:0, Y:0},
					zoom: p.Float64(12),
					angle: p.Float64(0),
					pitch: p.Float64(0),
				},
			},
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}
