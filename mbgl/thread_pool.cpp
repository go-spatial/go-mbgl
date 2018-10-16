#include <iostream>

#include <mbgl/util/default_thread_pool.hpp>
#include <mbgl/actor/scheduler.hpp>

#include "thread_pool.h"


MbglThreadPool * mbgl_thread_pool_new(int threads) {
    auto tpool = new mbgl::ThreadPool(threads);

    std::cout << "new ptr " << tpool << std::endl;

    return reinterpret_cast<MbglThreadPool*>(tpool);
}


void mbgl_thread_pool_destruct(MbglThreadPool * self) {
    auto cast = reinterpret_cast<mbgl::ThreadPool*>(self);
    std::cout << "old ptr " << self << std::endl;
    // TODO: GDEY/EAR7H need to fix this.
    //delete cast;
}
