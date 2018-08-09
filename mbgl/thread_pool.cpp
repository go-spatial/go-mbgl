#include <mbgl/util/default_thread_pool.hpp>

#include "thread_pool.h"

using namespace mbgl;

MbglThreadPool * mbgl_thread_pool_new(int threads) {
    auto tpool = new ThreadPool(threads);
    return reinterpret_cast<MbglThreadPool*>(tpool);
}


void mbgl_thread_pool_destruct(MbglThreadPool * ptr) {
    delete reinterpret_cast<ThreadPool*>(ptr);
}