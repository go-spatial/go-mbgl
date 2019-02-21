package bounds

import (
	"strconv"
	"testing"

	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/cmp"
)

func TestZoom(t *testing.T) {

	type tcase struct {
		bounds *geom.Extent
		width  float64
		height float64
		zoom   float64
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			zoom := Zoom(tc.bounds, tc.width, tc.height)
			if !cmp.Float(tc.zoom, zoom) {
				t.Errorf("zoom, expected %v got %v", tc.zoom, zoom)
			}
		}
	}

	tests := [...]tcase{
		{
			// {west, south, east, north}
			bounds: &geom.Extent{
				-117.1673735976219,  // west
				32.71965828903011,   // south
				-117.16439634561537, // east
				32.7204706651118,    // north
			},
			width:  862,
			height: 300,
			zoom:   18.0,
		},
	}

	for i := range tests {
		t.Run(strconv.Itoa(i), fn(tests[i]))
	}

}

func TestCenterZoom(t *testing.T) {
	type tcase struct {
		bounds *geom.Extent
		width  float64
		height float64
		zoom   float64
		center [2]float64
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {

			center, zoom := CenterZoom(tc.bounds, tc.width, tc.height)
			if !(cmp.Float(tc.center[0], center[0]) && cmp.Float(tc.center[0], center[0])) {
				t.Errorf("center, expected %v got %v", tc.center, center)
				return
			}

			if !cmp.Float(tc.zoom, zoom) {
				t.Errorf("zoom, expected %v got %v", tc.zoom, zoom)
				return
			}

		}
	}

	tests := [...]tcase{
		{
			// {west, south, east, north}
			bounds: &geom.Extent{
				-117.147086641189, // west
				32.7305263087481,  // south
				-117.180183060805, // east
				32.6963180459813,  // north
			},
			width:  1107,
			height: 360,
			zoom:   13.0,
			center: [2]float64{32.71342381720108, -117.163634850997},
		},
		{
			// {west, south, east, north}
			bounds: &geom.Extent{
				-117.1673735976219,  // west
				32.71965828903011,   // south
				-117.16439634561537, // east
				32.7204706651118,    // north
			},
			width:  1107,
			height: 360,
			zoom:   18.0,
			center: [2]float64{32.720064477996, -117.16588497161865},
		},
	}

	for i := range tests {
		t.Run(strconv.Itoa(i), fn(tests[i]))
	}

}

func TestTransform(t *testing.T) {
	type subcase struct {
		point [2]float64
		scale float64
		pt    [2]float64
	}

	type tcase struct {
		prj   AProjection
		cases []subcase
	}

	fn := func(tc tcase) (string, func(*testing.T)) {
		return tc.prj.String(), func(t *testing.T) {

			fn := func(prj AProjection, tc subcase) func(*testing.T) {
				return func(t *testing.T) {
					t.Run("transform", func(t *testing.T) {

						pt := prj.Transform(tc.point, tc.scale)
						if !(cmp.Float(pt[0], tc.pt[0]) && cmp.Float(pt[1], tc.pt[1])) {
							t.Errorf(" %v Transform, expected %v got %v", prj, tc.pt, pt)
							t.Logf("%v %v ", cmp.Float(pt[0], tc.pt[0]), cmp.Float(pt[1], tc.pt[1]))
						}
					})
					t.Run("untransform", func(t *testing.T) {

						point := prj.Untransform(tc.pt, tc.scale)
						if !(cmp.Float(point[0], tc.point[0]) && cmp.Float(point[1], tc.point[1])) {
							t.Errorf(" %v Transform, expected %v got %v", prj, tc.point, point)
							t.Logf("%v %v ", cmp.Float(point[0], tc.point[0]), cmp.Float(point[1], tc.point[1]))
						}
					})

				}
			}
			for i := range tc.cases {
				t.Run(strconv.Itoa(i), fn(tc.prj, tc.cases[i]))
			}
		}
	}

	tests := [...]tcase{
		{
			prj: ESPG3857,
			cases: []subcase{
				{
					point: [2]float64{44.68203449249269, 103.35370445251465},
					scale: 2.0,
					pt:    [2]float64{1.0000022299196951, 0.9999948419882011},
				},
				{
					point: [2]float64{1.2961, 103.831},
					pt:    [2]float64{0.5000000323418455, 0.49999740908404816},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fn(tc))
	}

}

func TestProject(t *testing.T) {
	type subcase struct {
		point [2]float64
		pt    [2]float64
	}

	type tcase struct {
		prj   AProjection
		cases []subcase
	}

	fn := func(tc tcase) (string, func(*testing.T)) {
		return tc.prj.String(), func(t *testing.T) {

			fn := func(prj AProjection, tc subcase) func(t *testing.T) {
				return func(t *testing.T) {

					t.Run("Project", func(t *testing.T) {

						pt := prj.Project(tc.point)
						if !(cmp.Float(pt[0], tc.pt[0]) && cmp.Float(pt[1], tc.pt[1])) {
							t.Errorf(" %v Project %v, expected %v got %v", prj, tc.point, tc.pt, pt)
							t.Logf("%v %v ", cmp.Float(pt[0], tc.pt[0]), cmp.Float(pt[1], tc.pt[1]))
						}

					})

					t.Run("Unproject", func(t *testing.T) {
						point := prj.Project(tc.pt)
						if !(cmp.Float(point[0], tc.point[0]) && cmp.Float(point[1], tc.point[1])) {
							t.Errorf(" %v Unproject %v, expected %v got %v", prj, tc.pt, tc.point, point)
							t.Logf("%v %v ", cmp.Float(point[0], tc.point[0]), cmp.Float(point[1], tc.point[1]))
						}

					})

				}
			}
			for i := range tc.cases {
				t.Run(strconv.Itoa(i), fn(tc.prj, tc.cases[i]))
			}
		}
	}
	tests := [...]tcase{
		{
			prj: ESPG3857,
			cases: []subcase{
				{
					point: [2]float64{32.7305263087481, -117.180183060805},
					pt:    [2]float64{-13044438.309391394, 3859590.2188198487},
				},
				{
					point: [2]float64{0.5000000323418455, 0.49999740908404816},
					pt:    [2]float64{55659.45697719234, 55660.45546583664},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(fn(tc))
	}
}
