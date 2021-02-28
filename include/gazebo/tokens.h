#ifndef TOKENS_H
#define TOKENS_H

#include <gazebo/string.h>

#define TOKEN_FOREACH(CB) \
    CB(TK_INVALID)        \
    CB(TK_COMMENT)        \
    CB(TK_WHITESPACE)     \
    CB(TK_SEMICOLON)      \
    CB(TK_DOT)            \
    CB(TK_COMMA)          \
    CB(TK_EQUAL)          \
    CB(TK_EQUAL_EQUAL)    \
    CB(TK_NOT)            \
    CB(TK_NOT_EQUAL)      \
    CB(TK_GREATER)        \
    CB(TK_GREATER_EQUAL)  \
    CB(TK_LESS)           \
    CB(TK_LESS_EQUAL)     \
    CB(TK_IF)             \
    CB(TK_ELSE)           \
    CB(TK_RETURN)         \
    CB(TK_WHILE)          \
    CB(TK_FOR)            \
    CB(TK_DEL)            \
    CB(TK_FUN)

#define _TOKEN_ENUM(name) name,
#define _TOKEN_NAME(name) #name,

typedef enum { TOKEN_FOREACH(_TOKEN_ENUM) } token_type;

static const char* token_type_names[] = {TOKEN_FOREACH(_TOKEN_NAME)};

static inline const char* token_type_name(token_type type)
{
    return token_type_names[type];
}

typedef struct {
    token_type type;
    char*      source;
    int        length;
} token;

static inline token token_new(token_type type, char* source, int length)
{
    return (token){
        .type   = type,
        .source = source,
        .length = length,
    };
}

#endif /* TOKENS_H */
