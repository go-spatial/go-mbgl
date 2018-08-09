package mbgl
/*
#include "scheduler.h"
#include "thread_pool.h"
 */
import "C"
import (
	"runtime"
)

type ThreadPool struct {
	ptr *C.MbglThreadPool
}

func NewThreadPool(threads int) (*ThreadPool) {
	ret := &ThreadPool{
		ptr: C.mbgl_thread_pool_new(C.int(threads)),
	}

	//fmt.Printf("new %p\n", ret)

	runtime.SetFinalizer(ret, func(pool *ThreadPool) {
		//fmt.Println("alloc finalizer")
		pool.Destruct()
	})

	//fmt.Println("return")

	return ret
}

func (t ThreadPool) threadPool () *C.MbglThreadPool {
	return t.ptr
}

// Scheduler is a prarent class
func (t ThreadPool) scheduler() *C.MbglScheduler {
	return (*C.MbglScheduler)(t.ptr)
}

func (t *ThreadPool) Destruct() {
	if t.ptr != nil {
		C.mbgl_thread_pool_destruct(t.ptr)
		t.ptr = nil
	}
}