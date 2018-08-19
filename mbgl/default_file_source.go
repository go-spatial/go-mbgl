package mbgl

/*
#include "file_source.h"
#include "default_file_source.h"
 */
import "C"

type DefaultFileSource C.MbglDefaultFileSource

func (fs *DefaultFileSource) defaultFileSource() *C.MbglDefaultFileSource {
	return (*C.MbglDefaultFileSource)(fs)
}

// This instantiates a new file source which can handle online and offline sources. It will create a chache file
func NewDefaultFileSource(cachePath, assetRoot string, maxCache *uint64) *DefaultFileSource {
	ptr :=  C.mbgl_default_file_source_new(C.CString(cachePath),
		C.CString(assetRoot),
		(*C.uint64_t)(maxCache))

	return (*DefaultFileSource)(ptr)
}

func (fs *DefaultFileSource) fileSource() *C.MbglFileSource {
	return (*C.MbglFileSource)(fs.defaultFileSource())
}

func (fs *DefaultFileSource) Destruct() {
	C.mbgl_default_file_source_destruct(fs.defaultFileSource())
}