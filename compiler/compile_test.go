package compiler

import (
	"testing"

	"github.com/johnfrankmorgan/gazebo/gvalue"
	"github.com/johnfrankmorgan/gazebo/op"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	assert := assert.New(t)

	code := Compile("(let x 5) (let y (+ x 5)) (print y)")

	exp := []op.Instruction{
		op.LoadConst.Instruction(gvalue.New(5)),
		op.StoreName.Instruction(gvalue.New("x")),
		op.LoadName.Instruction(gvalue.New("+")),
		op.LoadName.Instruction(gvalue.New("x")),
		op.LoadConst.Instruction(gvalue.New(5)),
		op.CallFunc.Instruction(gvalue.New(2)),
		op.StoreName.Instruction(gvalue.New("y")),
		op.LoadName.Instruction(gvalue.New("print")),
		op.LoadName.Instruction(gvalue.New("y")),
		op.CallFunc.Instruction(gvalue.New(1)),
	}

	assert.Equal(exp, code)
}
