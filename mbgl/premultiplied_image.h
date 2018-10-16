#ifndef _premultiplied_image_h
#define _premultiplied_image_h

#include <stdlib.h>
#include <stdint.h>

typedef struct{} MbglMapSnapshotter;
typedef struct{} MbglPremultipliedImage;

typedef struct{
    size_t height;
    size_t width;
    uint8_t * data;
} RawImage;

#ifdef __cplusplus
extern "C" {
#endif

// image
RawImage * mbgl_premultiplied_image_raw(MbglPremultipliedImage * self);

void mbgl_premultiplied_image_destruct(MbglPremultipliedImage * self);

#ifdef __cplusplus
} // extern "C"
#endif


#endif
