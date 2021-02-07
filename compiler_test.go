package gazebo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	assert := assert.New(t)

	code := Compile("(if (> 0 1) (println true) (println false))")

	code.Dump()

	exp := Code{
		OpLoadName.Ins(">"),
		OpLoadConst.Ins(float64(0)),
		OpLoadConst.Ins(float64(1)),
		OpCallFunc.Ins(2),
		OpRelJumpIfTrue.Ins(4),
		OpLoadName.Ins("println"),
		OpLoadName.Ins("false"),
		OpCallFunc.Ins(1),
		OpRelJump.Ins(3),
		OpLoadName.Ins("println"),
		OpLoadName.Ins("true"),
		OpCallFunc.Ins(1),
	}

	assert.Equal(exp, code)
}
