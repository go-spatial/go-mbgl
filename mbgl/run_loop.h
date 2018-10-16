#ifndef _mbgl_run_loop_h
#define _mbgl_run_loop_h

typedef struct {} MbglRunLoop;


#ifdef __cplusplus
extern "C" {
#endif

MbglRunLoop * mbgl_run_loop_new();
void mbgl_run_loop_destruct(MbglRunLoop * self);

#ifdef __cplusplus
} // extern "C"
#endif 

#endif