package compiler

import (
	"regexp"
	"strings"

	"github.com/johnfrankmorgan/gazebo/assert"
)

func subexpr(tokens []string, start int) []string {
	assert.True(tokens[start] == "(")

	depth := 0

	for idx, token := range tokens[start:] {
		if token == "(" {
			depth++
		} else if token == ")" {
			depth--
			if depth == 0 {
				return tokens[start+1 : start+idx]
			}
		}
	}

	assert.Unreached()
	return nil
}

func split(source string) []string {
	normalize := func(tokens []string) []string {
		parsed := make([]string, 0)

		for _, tk := range tokens {
			tk = strings.TrimSpace(tk)
			if tk != "" {
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

type sexp struct {
	children []*sexp
	value    string
}

func (m *sexp) isAtom() bool {
	return m.value != ""
}

type parser struct {
	tokens   []string
	position int
	sexp     *sexp
}

func (m *parser) finished() bool {
	return m.position >= len(m.tokens)
}

func (m *parser) next() string {
	token := m.tokens[m.position]
	m.position++
	return token
}

func (m *parser) parse() {
	m.sexp = new(sexp)

	for !m.finished() {
		token := m.next()

		switch token {
		case "(":
			subp := &parser{tokens: subexpr(m.tokens, m.position-1)}
			subp.parse()
			m.sexp.children = append(m.sexp.children, subp.sexp)
			m.position += subp.position + 1

		default:
			m.sexp.children = append(m.sexp.children, &sexp{value: token})
		}
	}
}
