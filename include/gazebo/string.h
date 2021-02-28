#ifndef STRING_H
#define STRING_H

typedef struct {
    char* value;
    int   length;
    int   capacity;
} g_string;

g_string* g_string_new(char*);

void g_string_free(g_string*);

void g_string_append(g_string*, g_string*);

void g_string_append_chars(g_string*, char*);

static inline char* g_string_value(g_string* str)
{
    return str->value;
}

static inline int g_string_length(g_string* str)
{
    return str->length;
}

static inline int g_string_capacity(g_string* str)
{
    return str->capacity;
}

#endif /* STRING_H */
