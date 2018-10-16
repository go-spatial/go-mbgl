package mbgl

/*
#include "run_loop.h"
*/
import "C"

type RunLoop C.MbglRunLoop

func NewRunLoop() *RunLoop {
	ptr := C.mbgl_run_loop_new()

	return (*RunLoop)(ptr)
}

func (rl *RunLoop) Destruct() {
	C.mbgl_run_loop_destruct((*C.MbglRunLoop)(rl))
}
