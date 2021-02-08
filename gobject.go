package gazebo

import (
	"fmt"
	"math"
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

func (m *GFuncCtx) Self() *GObject {
	assert.True(len(m.Args) > 0)
	return m.Args[0]
}

func (m *GFuncCtx) Parse(args ...interface{}) {
	assert.True(len(m.Args) >= len(args))

	for i, arg := range args {
		value := m.Args[i].Value

		switch arg := arg.(type) {
		case *bool:
			*arg = value.(bool)

		case *int:
			*arg = int(value.(float64))

		case *float64:
			*arg = value.(float64)

		case *string:
			*arg = value.(string)

		case **GObject:
			*arg = m.Args[i]

		default:
			assert.Unreached("cannot parse arg type %T", arg)
		}
	}
}

func (m *GFuncCtx) Interfaces() []interface{} {
	ifaces := make([]interface{}, len(m.Args))

	for i, arg := range m.Args {
		ifaces[i] = arg.Interface()
	}

	return ifaces
}

type GFunc func(*GFuncCtx) *GObject

type GMethods map[string]GFunc

type GType struct {
	Name    string
	Parent  *GType
	Methods GMethods
}

func (m *GType) Resolve(name string) GFunc {
	if method, ok := m.Methods[name]; ok {
		return method
	}

	if m.Parent != nil {
		return m.Parent.Resolve(name)
	}

	return nil
}

func (m *GType) Implements(name string) bool {
	return m.Resolve(name) != nil
}

var gtypes struct {
	Base   *GType
	Nil    *GType
	Bool   *GType
	Number *GType
	String *GType
	Func   *GType
}

func init() {
	gtypes.Base = &GType{
		Name:   "Base",
		Parent: nil,
		Methods: GMethods{
			"inspect": GFunc(func(ctx *GFuncCtx) *GObject {
				self := ctx.Self()

				inspection := fmt.Sprintf(
					"<gtypes.%s %p>(%v %p)",
					self.Type.Name,
					self.Type,
					self.Interface(),
					self,
				)

				return NewGObjectInferred(inspection)
			}),
		},
	}

	gtypes.Nil = &GType{
		Name:   "Nil",
		Parent: gtypes.Base,
	}

	gtypes.Bool = &GType{
		Name:   "Bool",
		Parent: gtypes.Base,
	}

	gtypes.Number = &GType{
		Name:   "Number",
		Parent: gtypes.Base,
	}

	gtypes.String = &GType{
		Name:   "String",
		Parent: gtypes.Base,
		Methods: GMethods{
			"replace": GFunc(func(ctx *GFuncCtx) *GObject {
				var (
					self    string
					search  string
					replace string
				)

				ctx.Parse(&self, &search, &replace)

				return NewGObjectInferred(strings.ReplaceAll(self, search, replace))
			}),
		},
	}

	gtypes.Func = &GType{
		Name:   "Func",
		Parent: gtypes.Base,
	}

	initbuiltins()
}

type GObject struct {
	Type  *GType
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

	}

	assert.Unreached("Could not infer type for %T %v", val, val)
	return nil
}

func (m *GObject) Interface() interface{} {
	if value, ok := m.Value.(float64); ok && math.Mod(value, 1) == 0 {
		return int(value)
	}

	return m.Value
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

func (m *GObject) Call(name string, ctx *GFuncCtx) *GObject {
	assert.True(m.Type.Implements(name))

	return m.Type.Resolve(name)(ctx)
}
