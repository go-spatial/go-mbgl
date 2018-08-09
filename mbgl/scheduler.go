package mbgl

/*
#include "scheduler.h"
 */
import "C"

// Represents the scheduler super class
type Scheduler interface {
	scheduler() *C.MbglScheduler
}

// Satisfies the schduler super class
type scheduler struct {
	ptr *C.MbglScheduler
}

func (s scheduler) scheduler() *C.MbglScheduler {
	return s.ptr
}

func SchedulerGetCurrent() Scheduler {
	ptr := C.mbgl_scheduler_get_current()
	return scheduler{
		ptr: ptr,
	}
}

func SchedulerSetCurrent(sched Scheduler) {
	C.mbgl_scheduler_set_current(sched.scheduler())
}