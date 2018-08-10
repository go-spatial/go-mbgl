package mbgl_test

import (
	"testing"
	"github.com/go-spatial/go-mbgl/mbgl"
)

func TestScheduler(t *testing.T) {
	type tcase struct {
		sched mbgl.Scheduler
	}

	fn := func(tc tcase, t *testing.T) {
		mbgl.SchedulerSetCurrent(tc.sched)

		sched := mbgl.SchedulerGetCurrent()

		if (sched == nil) != (tc.sched == nil) {
			t.Fatalf("incorrect value %v, expected %v", sched, tc.sched)
		}
	}

	testcases := map[string]tcase{
		"1" : {
			sched: mbgl.NewThreadPool(1),
		},
		"2" : {
			sched: mbgl.NewThreadPool(4),
		},
		"nil" : {
			sched: nil,
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}

}
