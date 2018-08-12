#ifndef _mbgl_default_file_source_h
#define _mbgl_default_file_source_h

#include <stdint.h>

typedef struct {} MbglDefaultFileSource;

#ifdef __cplusplus
extern "C"{
#endif

MbglDefaultFileSource * mbgl_default_file_source_new(const char * cachePath, const char * assetRoot, uint64_t * maxCacheSize);

void mbgl_default_file_source_destruct(MbglDefaultFileSource * self);

#ifdef __cplusplus
} // extern "C"{
#endif

#endif // _mbgl_default_file_source_h