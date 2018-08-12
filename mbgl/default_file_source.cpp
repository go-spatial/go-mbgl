#include <cstdint>

#include <mbgl/storage/default_file_source.hpp>

#include "default_file_source.h"

using namespace mbgl;

MbglDefaultFileSource * mbgl_default_file_source_new(const char * cachePath, const char * assetRoot, uint64_t * maxCacheSize) {
    auto _cachePath = std::string(cachePath);
    auto _assetRoot = std::string(assetRoot);

    DefaultFileSource * fs;

    if (maxCacheSize != nullptr) {
        fs = new DefaultFileSource(_cachePath, _assetRoot, *maxCacheSize);
    } else {
        fs = new DefaultFileSource(_cachePath, _assetRoot);
    }

    return reinterpret_cast<MbglDefaultFileSource*>(fs);
}

void mbgl_default_file_source_destruct(MbglDefaultFileSource * self) {
    auto cast = reinterpret_cast<DefaultFileSource*>(self);
    delete cast;
}