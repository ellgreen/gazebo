package gazebo

import (
	"fmt"

	"github.com/johnfrankmorgan/gazebo/assert"
)

type GFuncArgCtx struct {
	VM   interface{}
	Args []*GObject
}

func (m *GFuncArgCtx) Expects(argc int) {
	assert.Len(m.Args, argc, "expected %d arguments, got %d", argc, len(m.Args))
}

func (m *GFuncArgCtx) ExpectsAtLeast(argc int) {
	assert.True(len(m.Args) >= argc, "expected at least %d arguments, got %d", argc, len(m.Args))
}

func (m *GFuncArgCtx) Self() *GObject {
	return m.Args[0]
}

type GFunc func(*GFuncArgCtx) *GObject

type GMethods map[string]GFunc

type GType struct {
	Name    string
	Parent  *GType
	Methods GMethods
}

func (m *GType) resolve(name string) GFunc {
	if meth, ok := m.Methods[name]; ok {
		return meth
	}

	if m.Parent != nil {
		return m.Parent.resolve(name)
	}

	return nil
}

func (m *GType) Implements(name string) bool {
	return m.resolve(name) != nil
}

var gtypemethods = struct {
	toBool   string
	toString string
	toNumber string
	isNil    string
}{
	toBool:   "?",
	toString: "str",
	toNumber: "num",
	isNil:    "nil?",
}

var gtypes = struct {
	init   func()
	Object *GType
	Nil    *GType
	Bool   *GType
	Number *GType
	String *GType
}{
	init: func() {
		gtypes.Object = &GType{
			Name: "Object",
			Methods: GMethods{
				gtypemethods.toBool: func(args *GFuncArgCtx) *GObject {
					return NewGObjectInferred(true)
				},

				gtypemethods.toString: func(args *GFuncArgCtx) *GObject {
					args.Expects(1)
					self := args.Self()
					str := fmt.Sprintf(
						"<gtypes.%s %p>(%p %v)",
						self.Type.Name,
						self.Type,
						self,
						self.Value,
					)
					return NewGObjectInferred(str)
				},

				gtypemethods.toNumber: func(args *GFuncArgCtx) *GObject {
					assert.Unreached("gtypemethods.toNumber not implemented")
					return nil
				},

				gtypemethods.isNil: func(args *GFuncArgCtx) *GObject {
					return NewGObjectInferred(false)
				},
			},
		}

		gtypes.Nil = &GType{
			Name: "Nil",
			Methods: GMethods{
				gtypemethods.toBool: func(args *GFuncArgCtx) *GObject {
					return NewGObjectInferred(false)
				},

				gtypemethods.isNil: func(args *GFuncArgCtx) *GObject {
					return NewGObjectInferred(true)
				},
			},
		}

		gtypes.Bool = &GType{
			Name: "Bool",
			Methods: GMethods{
				gtypemethods.toBool: func(args *GFuncArgCtx) *GObject {
					args.Expects(1)
					return NewGObjectInferred(args.Self().Value.(bool))
				},
			},
		}

		gtypes.Number = &GType{
			Name:    "Number",
			Methods: GMethods{},
		}

		gtypes.String = &GType{
			Name:    "String",
			Methods: GMethods{},
		}
	},
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
		return &GObject{Type: gtypes.String, Value: val}
	}

	assert.Unreached("Could not infer type for %v", val)
	return nil
}

func (m *GObject) Call(name string, args *GFuncArgCtx) *GObject {
	method := m.Type.resolve(name)
	assert.NotNil(method)
	return method(args)
}
