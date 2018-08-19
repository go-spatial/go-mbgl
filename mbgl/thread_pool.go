package mbgl

/*
#include "scheduler.h"
#include "thread_pool.h"
 */
import "C"

type ThreadPool C.MbglThreadPool

func NewThreadPool(threads int) (*ThreadPool) {
	ptr := C.mbgl_thread_pool_new(C.int(threads))

	return (*ThreadPool)(ptr)
}

func (t *ThreadPool) threadPool() *C.MbglThreadPool {
	return (*C.MbglThreadPool)(t)
}

// Scheduler is a prarent class
func (t *ThreadPool) scheduler() *C.MbglScheduler {
	return (*C.MbglScheduler)(t)
}

func (t *ThreadPool) Destruct() {
	C.mbgl_thread_pool_destruct(t.threadPool())
}
