#ifndef _mbgl_thread_pool_h
#define _mbgl_thread_pool_h

typedef struct{} MbglThreadPool;

#ifdef __cplusplus
extern "C"{
#endif //__cplusplus

MbglThreadPool *mbgl_thread_pool_new(int threads);

void mbgl_thread_pool_destruct(MbglThreadPool * self);

#ifdef __cplusplus
}
#endif //__cplusplus

#endif //_mbgl_thread_pool_h