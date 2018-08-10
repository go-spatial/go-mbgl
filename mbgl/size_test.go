package mbgl

import "testing"

func TestSize(t *testing.T) {
	type size struct{
		width, height uint32
}

	type tcase struct {
		vals []size
	}

	fn := func(tc tcase, t *testing.T) {

		s := Size{}
		for _, v := range tc.vals {
			s.Height = v.height
			s.Width = v.width

			s.cSize()
			s.Destruct()

		}

		// run on destructed Size
		s.cSize()
		s.Destruct()

		// finalizer
		s = Size{}
	}

	testcases := map[string]tcase {
		"1": {
			[]size {
			 	{0, 0},
			},
		},
		"2": {
			[]size{
				{0, 0},
				{128, 128},
				{128, 1024},
			},
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}
