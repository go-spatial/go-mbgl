package mbgl

/*
#include "file_source.h"
 */
import "C"

type FileSource interface {
	fileSource() *C.MbglFileSource
	Destruct()
}
