package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	assert := assert.New(t)

	source := `
		; this is a comment
		(if (= 1.1 2 "test"))
	`

	expected := []tokentype{
		tkcomment,
		tkwhitespace,
		tkparenopen,
		tkident,
		tkwhitespace,
		tkparenopen,
		tkident,
		tkwhitespace,
		tknumber,
		tkwhitespace,
		tknumber,
		tkwhitespace,
		tkstring,
		tkparenclose,
		tkparenclose,
	}

	got := []tokentype{}
	for _, token := range tokenize(source) {
		got = append(got, token.typ)
	}

	assert.Equal(expected, got)
}
