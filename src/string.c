#include <string.h>

#include <gazebo.h>
#include <gazebo/string.h>

#ifdef G_DEBUG_STRING
#define G_DEBUG(...) _G_DEBUG(__VA_ARGS__)
#else
#define G_DEBUG(...)
#endif

g_string* g_string_new(char* src)
{
    G_ASSERT(src);

    int       len = strlen(src);
    g_string* str = G_NEW(g_string);

    str->length   = len;
    str->capacity = len;
    str->value    = G_MALLOC(str->capacity + 1);

    strcpy(str->value, src);

    return str;
}

void g_string_free(g_string* str)
{
    G_ASSERT(str);

    G_FREE(str->value);
    G_FREE(str);
}

static inline void g_string_ensure_capacity(g_string* str, size_t length)
{
    G_ASSERT(str);

    if (length < (size_t)str->capacity) {
        return;
    }

    str->capacity = length + 1;
    str->value    = G_REALLOC(str->value, str->capacity);
}

void g_string_append(g_string* str, g_string* src)
{
    G_ASSERT(str && src);

    g_string_ensure_capacity(str, str->length + src->length);
    strcpy(str->value + str->length, src->value);
    str->length += src->length;
}

void g_string_append_chars(g_string* str, char* src)
{
    G_ASSERT(str && src);

    int length = strlen(src);

    g_string_ensure_capacity(str, str->length + length);
    strcpy(str->value + str->length, src);
    str->length += length;
}
