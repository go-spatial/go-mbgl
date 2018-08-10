#include <mbgl/actor/scheduler.hpp>

#include "scheduler.h"

using namespace mbgl;

// instance methods

void mbgl_scheduler_destruct(MbglScheduler * self) {
    auto cast = reinterpret_cast<Scheduler*>(self);
    delete self;
}

//static methods

MbglScheduler * mbgl_scheduler_get_current() {
    auto sched = Scheduler::GetCurrent();
    return reinterpret_cast<MbglScheduler*>(sched);
}

void mbgl_scheduler_set_current(MbglScheduler * sched) {
    auto cast = reinterpret_cast<Scheduler*>(sched);
    Scheduler::SetCurrent(cast);
}

