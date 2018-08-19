package mbgl

/*
#include "scheduler.h"
 */
import "C"

// Represents the scheduler super class
type Scheduler interface {
	scheduler() *C.MbglScheduler
	Destruct()
}

// Satisfies the schduler inteface
type scheduler C.MbglScheduler

func (s *scheduler) scheduler() *C.MbglScheduler {
	return (*C.MbglScheduler)(s)
}

func (s *scheduler) Destruct() {
	C.mbgl_scheduler_destruct(s.scheduler())
}

func SchedulerGetCurrent() Scheduler {
	ptr := C.mbgl_scheduler_get_current()
	if ptr == nil {
		return nil
	}

	return (*scheduler)(ptr)
}

func SchedulerSetCurrent(sched Scheduler) {
	if sched == nil {
		C.mbgl_scheduler_set_current(nil)
	} else {
		C.mbgl_scheduler_set_current(sched.scheduler())
	}
}