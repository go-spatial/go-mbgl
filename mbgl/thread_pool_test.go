package mbgl_test

import (
	"testing"
	"runtime"
	"time"
	"github.com/go-spatial/go-mbgl/mbgl"
)

func TestThreadPool(t *testing.T) {
	tpool := mbgl.NewThreadPool(4)
	tpool.Destruct()
}

func TestThreadPoolFinalizer(t *testing.T) {

	c := make(chan struct{})

	f := func() {
		tpool := mbgl.NewThreadPool(4)
		// fmt.Printf("pointer: %p\n", tpool)

		// clear previous
		// TODO (@ear7h): two finalizers end up running during this test, the one set in `New` and the test finalizer. this shouldn't happen.
		runtime.SetFinalizer(tpool, nil)
		runtime.SetFinalizer(tpool, func(pool *mbgl.ThreadPool) {
			pool.Destruct()
			c <- struct{}{}
		})

		// fmt.Println("test finalizer set: ", tpool.ptr)

	}

	f()

	time.Sleep(time.Second)
	runtime.GC()
	select {
	case <-c:
		//pass
	case <-time.After(time.Second):
		t.Fatalf("finalizer timeout")
	}
}