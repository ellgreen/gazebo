package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserParse(t *testing.T) {
	assert := assert.New(t)

	p := &parser{tokens: split("(if true (1 2 3) (4 5 6))")}
	p.parse()

	assert.Len(p.sexp.children, 1)
	assert.Len(p.sexp.children[0].children, 4)
	assert.Len(p.sexp.children[0].children[2].children, 3)
	assert.Len(p.sexp.children[0].children[3].children, 3)

	assert.True(p.sexp.children[0].children[2].children[0].isAtom())
}
