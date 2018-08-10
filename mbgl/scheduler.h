#ifndef _mbgl_scheduler_h
#define _mbgl_scheduler_h

typedef struct{} MbglScheduler;

#ifdef __cplusplus
extern "C"{
#endif

// instance methods
void mbgl_scheduler_destruct(MbglScheduler * self);

// Static mehtods
MbglScheduler * mbgl_scheduler_get_current();
void mbgl_scheduler_set_current(MbglScheduler *);


#ifdef __cplusplus
} // extern "C"
#endif

#endif //_mbgl_scheduler_h