#include <gazebo.h>
#include <gazebo/lexer.h>
#include <gazebo/tokens.h>

lexer* lexer_from_file(char* path)
{
    G_ASSERT(path);

    lexer* l       = G_NEW(lexer);
    l->source      = g_string_new("");
    l->position    = 0;
    l->token_start = 0;

    G_ASSERT(fs_read_file(path, l->source));

    return l;
}

void lexer_free(lexer* lexer)
{
    G_ASSERT(lexer);

    G_FREE(lexer->source);
    G_FREE(lexer);
}

static inline int finished(lexer* l)
{
    return l->position >= g_string_length(l->source);
}

static inline char peek(lexer* l)
{
    if (finished(l)) {
        return EOF;
    }

    return g_string_value(l->source)[l->position];
}

static inline char next(lexer* l)
{
    char ch = peek(l);

    if (ch != EOF) {
        ++l->position;
    }

    return ch;
}

static inline token new_token(lexer* l, token_type type)
{
    char* source   = g_string_value(l->source) + l->token_start;
    token tk       = token_new(type, source, l->position - l->token_start);
    l->token_start = l->position;
    return tk;
}

static inline int check(lexer* l, char ch)
{
    return peek(l) == ch;
}

static inline int match(lexer* l, char ch)
{
    if (check(l, ch)) {
        next(l);
        return 1;
    }

    return 0;
}

static inline token if_match(lexer* l, char ch, token_type type, token_type fallback)
{
    if (match(l, ch)) {
        return new_token(l, type);
    }

    return new_token(l, fallback);
}

token lexer_lex(lexer* l)
{
    G_ASSERT(l);

    char ch = next(l);

    switch (ch) {
    // comment
    // whitespace
    case ';':
        return new_token(l, TK_SEMICOLON);

    case '.':
        return new_token(l, TK_DOT);

    case ',':
        return new_token(l, TK_COMMA);

    case '=':
        return if_match(l, '=', TK_EQUAL_EQUAL, TK_EQUAL);
    case '!':
        return if_match(l, '=', TK_NOT_EQUAL, TK_NOT);

    case '>':
        return if_match(l, '=', TK_GREATER_EQUAL, TK_GREATER);

    case '<':
        return if_match(l, '=', TK_LESS_EQUAL, TK_LESS);
    }

    G_UNREACHED();
}
