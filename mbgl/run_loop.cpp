#include <mbgl/util/run_loop.hpp>

#include "run_loop.h"

using namespace mbgl::util;

MbglRunLoop * mbgl_run_loop_new() {
	auto ptr = new RunLoop();

	return reinterpret_cast<MbglRunLoop *>(ptr);
}

void mbgl_run_loop_destruct(MbglRunLoop * self) {
	auto _self = reinterpret_cast<RunLoop *>(self);

	delete _self;
}