#include <gazebo.h>
#include <gazebo/lexer.h>

int main()
{
    lexer* l = lexer_from_file("tests/test.gaz");

    while (1) {
        token tk = lexer_lex(l);
        printf("%s :: %.*s\n", token_type_name(tk.type), tk.length, tk.source);
    }

    lexer_free(l);

    return 0;
}
