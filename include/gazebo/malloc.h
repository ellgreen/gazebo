#ifndef MALLOC_H
#define MALLOC_H

#include <stdlib.h>

#define G_NEW(type) G_MALLOC(sizeof(type))
#define G_MALLOC(size) g_malloc(size)
#define G_REALLOC(ptr, size) g_realloc(ptr, size)
#define G_FREE(ptr) g_free(ptr)

void* g_malloc(size_t);

void* g_realloc(void*, size_t);

void g_free(void*);

#endif /* MALLOC_H */
