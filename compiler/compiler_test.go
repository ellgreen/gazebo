package compiler

import (
	"testing"

	"github.com/johnfrankmorgan/gazebo/compiler/op"
	"github.com/johnfrankmorgan/gazebo/debug"
	"github.com/johnfrankmorgan/gazebo/errors"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	debug.Enable()
	defer debug.Disable()

	assert := assert.New(t)

	code, err := Compile("(if (> 0 1) (println true) (println false)) (println i[1] t[.test])")

	assert.Nil(err)

	exp := Code{
		op.LoadName.Ins(">"),
		op.LoadConst.Ins(float64(0)),
		op.LoadConst.Ins(float64(1)),
		op.CallFunc.Ins(2),
		op.RelJumpIfTrue.Ins(4),
		op.LoadName.Ins("println"),
		op.LoadName.Ins("false"),
		op.CallFunc.Ins(1),
		op.RelJump.Ins(3),
		op.LoadName.Ins("println"),
		op.LoadName.Ins("true"),
		op.CallFunc.Ins(1),
		op.LoadName.Ins("println"),
		op.LoadName.Ins("i"),
		op.LoadConst.Ins(float64(1)),
		op.IndexGet.Ins(nil),
		op.LoadName.Ins("t"),
		op.AttributeGet.Ins("test"),
		op.CallFunc.Ins(2),
	}

	assert.Equal(exp, code)
}

func TestCompileError(t *testing.T) {
	assert := assert.New(t)

	_, err := Compile("(")
	assert.ErrorIs(err, errors.ErrEOF)
}
