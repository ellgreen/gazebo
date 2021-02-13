package compiler

import (
	"strings"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/debug"
)

func parse(source string) *sexpr {
	tokens := tokenize(source)

	if debug.Enabled() {
		tokens.dump()
	}

	parser := parser{tokens: tokens}

	return parser.parse()
}

type sexpr struct {
	children []*sexpr
	token    token
}

func (m *sexpr) dump(depth int) {
	indent := strings.Repeat(" ", depth*2)

	if m.atom() {
		debug.Printf("%s%s(%v)\n", indent, m.token.typ.name(), m.token.value)
		return
	}

	debug.Printf("%s(\n", indent)

	for _, expr := range m.children {
		expr.dump(depth + 1)
	}

	debug.Printf("%s)\n", indent)
}

func (m *sexpr) atom() bool {
	return m.token.typ.valid()
}

type parser struct {
	tokens   tokens
	position int
}

func (m *parser) unexpectedeof() *sexpr {
	assert.Unreached("unexpected eof at token offset %d", m.position)
	return nil
}

func (m *parser) finished() bool {
	return m.position >= len(m.tokens)
}

func (m *parser) peek() token {
	return m.tokens[m.position]
}

func (m *parser) next() token {
	token := m.tokens[m.position]
	m.position++
	return token
}

func (m *parser) subexpr(start int) []token {
	assert.True(m.tokens[start].is(tkparenopen))

	depth := 0

	for idx, token := range m.tokens[start:] {
		if token.is(tkparenopen) {
			depth++
		} else if token.is(tkparenclose) {
			depth--
			if depth == 0 {
				return m.tokens[start+1 : start+idx]
			}
		}
	}

	assert.Unreached("unterminated expression near token offset: %d", start)
	return nil
}

func (m *parser) parse() *sexpr {
	expr := new(sexpr)

	for !m.finished() {
		token := m.next()

		if token.is(tkwhitespace, tkcomment) {
			continue
		}

		if !token.is(tkparenopen) {
			expr.children = append(expr.children, &sexpr{
				token: token,
			})

			continue
		}

		parser := parser{tokens: m.subexpr(m.position - 1)}
		subexpr := parser.parse()
		expr.children = append(expr.children, subexpr)

		m.position += parser.position + 1
	}

	return expr
}
