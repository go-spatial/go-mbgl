package mbgl

import (
	"testing"
	"fmt"
)

func TestScheduler(t *testing.T) {
	type tcase struct {
		sched Scheduler
	}

	fn := func(tc tcase, t *testing.T) {
		SchedulerSetCurrent(tc.sched)

		sched := SchedulerGetCurrent()

		if (sched == nil) != (tc.sched == nil) {
			fmt.Printf("%p, %v\n", tc.sched, tc.sched == nil)
			t.Fatalf("incorrect value %v, expected %v", sched, tc.sched)
		} else if tc.sched == nil {
			// they are both nil, pass
		} else if sched.scheduler() != tc.sched.scheduler() {
			t.Fatalf("incorrect value %v, expected %v",sched.scheduler(), tc.sched.scheduler())
		}
	}

	testcases := map[string]tcase{
		"1" : {
			sched: NewThreadPool(1),
		},
		"2" : {
			sched: NewThreadPool(4),
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