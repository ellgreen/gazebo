package gazebo

import (
	"strings"

	"github.com/johnfrankmorgan/gazebo/assert"
)

type GFuncCtx struct {
	VM   *VM
	Args []*GObject
}

func (m *GFuncCtx) Expects(argc int) {
	assert.Len(m.Args, argc, "expected %d arguments, got %d", argc, len(m.Args))
}

func (m *GFuncCtx) ExpectsAtLeast(argc int) {
	assert.True(len(m.Args) >= argc, "expected at least %d arguments, got %d", argc, len(m.Args))
}

type GFunc func(*GFuncCtx) *GObject

type GType int

var gtypes = struct {
	Nil    GType
	Bool   GType
	Number GType
	String GType
	Func   GType
}{
	Nil:    1,
	Bool:   2,
	Number: 3,
	String: 4,
	Func:   5,
}

type GObject struct {
	Type  GType
	Value interface{}
}

func NewGObjectInferred(val interface{}) *GObject {
	switch val := val.(type) {
	case nil:
		return &GObject{Type: gtypes.Nil}

	case bool:
		return &GObject{Type: gtypes.Bool, Value: val}

	case int:
		return &GObject{Type: gtypes.Number, Value: float64(val)}

	case float64:
		return &GObject{Type: gtypes.Number, Value: val}

	case string:
		val = strings.ReplaceAll(val, "\\n", "\n")
		return &GObject{Type: gtypes.String, Value: val}

	case func(*GFuncCtx) *GObject:
		return &GObject{Type: gtypes.Func, Value: GFunc(val)}

	case GFunc:
		return &GObject{Type: gtypes.Func, Value: val}

	}

	assert.Unreached("Could not infer type for %v", val)
	return nil
}

func (m *GObject) IsTruthy() bool {
	switch m.Type {
	case gtypes.Nil:
		return false

	case gtypes.Bool:
		return m.Value.(bool)

	case gtypes.Number:
		return m.Value.(float64) != 0

	case gtypes.String:
		return m.Value.(string) != ""

	case gtypes.Func:
		return true
	}

	assert.Unreached("unknown type: %d", m.Type)
	return false
}
