#ifndef LEXER_H
#define LEXER_H

#include <gazebo.h>
#include <gazebo/tokens.h>

typedef struct {
    g_string* source;
    int       position;
    int       token_start;
} lexer;

lexer* lexer_from_file(char*);
void   lexer_free(lexer*);
token  lexer_lex(lexer*);

#endif /* LEXER_H */
