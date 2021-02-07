package compiler

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/gvalue"
	"github.com/johnfrankmorgan/gazebo/op"
)

var regexes = struct {
	numbers *regexp.Regexp
	strings *regexp.Regexp
	idents  *regexp.Regexp
}{
	numbers: regexp.MustCompile(`^-?[0-9]+(.[0-9]+)?$`),
	strings: regexp.MustCompile(`^".*"$`),
	idents:  regexp.MustCompile(`^[a-zA-Z!@$%^&*\/?<>_=+~-]+$`),
}

func atom(value string) []op.Instruction {
	switch true {
	case regexes.numbers.MatchString(value):
		value, err := strconv.ParseFloat(value, 64)
		assert.Nil(err)

		return []op.Instruction{op.LoadConst.Instruction(gvalue.New(value))}

	case regexes.strings.MatchString(value):
		value = value[1 : len(value)-1]
		value = strings.ReplaceAll(value, "\\n", "\n")
		return []op.Instruction{op.LoadConst.Instruction(gvalue.New(value))}

	case regexes.idents.MatchString(value):
		return []op.Instruction{op.LoadName.Instruction(gvalue.New(value))}
	}

	assert.Unreached("invalid atom: %q", value)
	return nil
}

func compile(exp *sexp) []op.Instruction {
	if exp.isAtom() {
		return atom(exp.value)
	}

	assert.True(len(exp.children) > 0)

	switch exp.children[0].value {
	case "let":
		assert.True(len(exp.children) == 3, "%# v", exp.children)
		assert.True(exp.children[1].isAtom())

		return append(
			compile(exp.children[2]),
			op.StoreName.Instruction(gvalue.New(exp.children[1].value)),
		)

	case "if":
		assert.True(len(exp.children) == 4, "%# v", exp.children)

		condition := exp.children[1]

		truePath := compile(exp.children[2])

		falsePath := append(
			compile(exp.children[3]),
			op.RelJump.Instruction(gvalue.New(len(truePath))),
		)

		code := append(
			compile(condition),
			op.RelJumpIfTrue.Instruction(gvalue.New(len(falsePath))),
		)

		code = append(code, falsePath...)

		return append(code, truePath...)

	case "while":
		assert.True(len(exp.children) == 3, "%# v", exp.children)

		body := compile(exp.children[2])

		condition := append(
			compile(exp.children[1]),
			op.RelJumpIfFalse.Instruction(gvalue.New(len(body)+1)),
		)

		body = append(
			body,
			op.RelJump.Instruction(gvalue.New(-len(body)-len(condition)-1)),
		)

		return append(condition, body...)
	}

	if !exp.children[0].isAtom() {
		code := make([]op.Instruction, 0)

		for _, exp := range exp.children {
			code = append(code, compile(exp)...)
		}

		return code
	}

	fun := compile(exp.children[0])
	argc := 0

	for _, arg := range exp.children[1:] {
		fun = append(fun, compile(arg)...)
		argc++
	}

	return append(fun, op.CallFunc.Instruction(gvalue.New(argc)))
}

func Compile(source string) []op.Instruction {
	var instructions []op.Instruction

	parser := parser{tokens: split(source)}
	parser.parse()

	for _, sexp := range parser.sexp.children {
		instructions = append(instructions, compile(sexp)...)
	}

	return instructions
}
