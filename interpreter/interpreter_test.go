package interpreter

import (
	"testing"

	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/gvalue"
	"github.com/stretchr/testify/assert"
)

func TestInterpreterEval(t *testing.T) {
	assert := assert.New(t)

	Create()
	code := compiler.Compile("(+ 4 16)")

	result := The().Eval(code)

	assert.Equal(gvalue.New(20), result)
}
