#include <mbgl/util/size.hpp>

#include "size.h"

using namespace mbgl;

MbglSize * mbgl_size_new(uint32_t width, uint32_t height) {
    auto s = new Size(width, height);

    return reinterpret_cast<MbglSize*>(s);
}

void mbgl_size_set(MbglSize * self, uint32_t width, uint32_t height) {
    auto obj = reinterpret_cast<Size*>(self);
    obj->width = width;
    obj->height = height;
}

void mbgl_size_destruct(MbglSize* self) {
    auto cast = reinterpret_cast<Size*>(self);
    delete cast;
}