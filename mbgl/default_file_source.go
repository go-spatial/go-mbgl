package mbgl

/*
#include "file_source.h"
#include "default_file_source.h"
 */
import "C"

type DefaultFileSource struct {
	ptr *C.MbglDefaultFileSource
}

func NewDefaultFileSource(cachePath, assetRoot string, maxCache *uint64) *DefaultFileSource {
	ptr :=  C.mbgl_default_file_source_new(C.CString(cachePath),
		C.CString(assetRoot),
		(*C.uint64_t)(maxCache))

	return &DefaultFileSource{
		ptr: ptr,
	}
}

func (fs DefaultFileSource) fileSource() *C.MbglFileSource {
	return (*C.MbglFileSource)(fs.ptr)
}

func (fs *DefaultFileSource) Destruct() {
	C.mbgl_default_file_source_destruct(fs.ptr)
	fs.ptr = nil
}