package mbgl

/*
#include "size.h"
 */
import "C"

type Size struct {
	Width, Height uint32
	ptr *C.MbglSize
}

func (s *Size) update() {

	if s.ptr == nil {
		s.ptr = C.mbgl_size_new(C.uint32_t(s.Width), C.uint32_t(s.Height))
	} else {
		C.mbgl_size_set(s.ptr, C.uint32_t(s.Width), C.uint32_t(s.Height))
	}
}

// called within the package
func (s *Size) cSize() *C.MbglSize {
	s.update()
	return s.ptr
}

func (s *Size) Destruct() {
	if s.ptr != nil {
		C.mbgl_size_destruct(s.ptr)
		s.ptr = nil
	}
}
