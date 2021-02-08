package gazebo

import (
	"regexp"
	"strings"

	"github.com/johnfrankmorgan/gazebo/assert"
)

type sexpr struct {
	children []*sexpr
	value    string
}

func (m *sexpr) isAtom() bool {
	return len(m.value) > 0
}

type parser struct {
	tokens   []string
	position int
}

func (m *parser) finished() bool {
	return m.position >= len(m.tokens)
}

func (m *parser) next() string {
	token := m.tokens[m.position]
	m.position++
	return token
}

func (m *parser) tokenize(source string) []string {
	normalize := func(tokens []string) []string {
		parsed := make([]string, 0)

		for _, tk := range tokens {
			if tk = strings.TrimSpace(tk); tk != "" {
				parsed = append(parsed, tk)
			}
		}

		return parsed
	}

	source = strings.ReplaceAll(source, "\n", " ")
	source = strings.ReplaceAll(source, "(", "( ")
	source = strings.ReplaceAll(source, ")", " )")

	tokens := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`).FindAllString(source, -1)

	return normalize(tokens)

}

func (m *parser) subexpr(start int) []string {
	assert.True(m.tokens[start] == "(")

	depth := 0

	for idx, token := range m.tokens[start:] {
		if token == "(" {
			depth++
		} else if token == ")" {
			depth--
			if depth == 0 {
				return m.tokens[start+1 : start+idx]
			}
		}
	}

	assert.Unreached()
	return nil
}

func (m *parser) parse(tokens []string) *sexpr {
	m.tokens = tokens
	m.position = 0

	expr := new(sexpr)

	for !m.finished() {
		token := m.next()

		if token != "(" {
			expr.children = append(expr.children, &sexpr{
				value: token,
			})

			continue
		}

		parser := &parser{}
		parsed := parser.parse(m.subexpr(m.position - 1))
		expr.children = append(expr.children, parsed)
		m.position += parser.position + 1
	}

	return expr
}
