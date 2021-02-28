#include <gazebo.h>
#include <gazebo/malloc.h>

#ifdef G_DEBUG_MEM
#define G_DEBUG(...) _G_DEBUG(__VA_ARGS__)
#else
#define G_DEBUG(...)
#endif

void* g_malloc(size_t size)
{
    void* ptr;

    if ((ptr = malloc(size))) {
        G_DEBUG("allocated %zu new bytes: %p", size, ptr);

        return ptr;
    }

    G_UNREACHED();
}

void* g_realloc(void* ptr, size_t size)
{
    void* new_ptr;

    if ((new_ptr = realloc(ptr, size))) {
        G_DEBUG("reallocated %zu new bytes: %p -> %p", size, ptr, new_ptr);

        return new_ptr;
    }

    G_UNREACHED();
}

void g_free(void* ptr)
{
    G_DEBUG("freeing %p", ptr);
    free(ptr);
}
