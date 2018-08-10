#ifndef _mbgl_size_h
#define _mbgl_size_h

#include <stdint.h>

typedef struct{} MbglSize;

#ifdef __cplusplus
extern "C"{
#endif //__cplusplus

MbglSize * mbgl_size_new(uint32_t width, uint32_t height);

void mbgl_size_set(MbglSize * self, uint32_t width, uint32_t height);

void mbgl_size_destruct(MbglSize * self);


#ifdef __cplusplus
}
#endif //__cplusplus

#endif // _mbgl_size_h