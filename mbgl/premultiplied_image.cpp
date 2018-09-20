
#include "premultiplied_image.h"
#include <mbgl/util/image.hpp>

// image

RawImage * mbgl_premultiplied_image_raw(MbglPremultipliedImage * img) {
    auto _img = reinterpret_cast<mbgl::PremultipliedImage*>(img);

    auto ret = new RawImage{};
    ret->height = _img->size.height;
    ret->width = _img->size.width;
    ret->data = _img->data.get();

    return ret;
}

void mbgl_premultiplied_image_destruct(MbglPremultipliedImage * self) {
    auto cast = reinterpret_cast<mbgl::PremultipliedImage*>(self);
    delete cast;
}
