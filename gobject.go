package gazebo

import (
	"fmt"
	"math"
	"strings"

	"github.com/johnfrankmorgan/gazebo/assert"
)

// GFuncCtx is used to pass arguments to builtin gazebo functions
type GFuncCtx struct {
	VM   *VM
	Args []*GObject
}

// Expects asserts that the number of arguments expects the provided value
func (m *GFuncCtx) Expects(argc int) {
	assert.Len(m.Args, argc, "expected %d arguments, got %d", argc, len(m.Args))
}

// ExpectsAtLeast asserts that at least X arguments have been specified
func (m *GFuncCtx) ExpectsAtLeast(argc int) {
	assert.True(len(m.Args) >= argc, "expected at least %d arguments, got %d", argc, len(m.Args))
}

// Self returns the function's first argument
func (m *GFuncCtx) Self() *GObject {
	assert.True(len(m.Args) > 0, "expected self argument, got 0 arguments")
	return m.Args[0]
}

// Parse parses arguments into the provided pointers
func (m *GFuncCtx) Parse(args ...interface{}) {
	assert.True(len(m.Args) >= len(args), "too few arguments to parse")

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

// Interfaces returns the interface values of the provided *GObjects
func (m *GFuncCtx) Interfaces() []interface{} {
	ifaces := make([]interface{}, len(m.Args))

	for i, arg := range m.Args {
		ifaces[i] = arg.Interface()
	}

	return ifaces
}

// GFunc is the type of gazebo's builtin functions
type GFunc func(*GFuncCtx) *GObject

// GUserFunc is the type of values store in a gtypes.UserFunc instance
type GUserFunc struct {
	params []string
	body   Code
	env    *env
}

// GMethods is a map of names to GFuncs
type GMethods map[string]GFunc

// GType represents the type of gazebo values
type GType struct {
	Name    string
	Parent  *GType
	Methods GMethods
}

// Resolve resolves a method on a GType
func (m *GType) Resolve(name string) GFunc {
	if method, ok := m.Methods[name]; ok {
		return method
	}

	if m.Parent != nil {
		return m.Parent.Resolve(name)
	}

	return nil
}

// Implements checks if a method is implemented by a GType
func (m *GType) Implements(name string) bool {
	return m.Resolve(name) != nil
}

var gtypes struct {
	Base     *GType
	Nil      *GType
	Bool     *GType
	Number   *GType
	String   *GType
	Func     *GType
	UserFunc *GType
	Internal *GType
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

	gtypes.UserFunc = &GType{
		Name:   "UserFunc",
		Parent: gtypes.UserFunc,
	}

	gtypes.Internal = &GType{
		Name:   "Internal",
		Parent: gtypes.Base,
	}

	initbuiltins()
}

// GObject is the type of all values in gazebo
type GObject struct {
	Type  *GType
	Value interface{}
}

// NewGObjectInferred creates an appropriate *GObject for the provided value
func NewGObjectInferred(value interface{}) *GObject {
	switch value := value.(type) {
	case nil:
		return &GObject{Type: gtypes.Nil}

	case bool:
		return &GObject{Type: gtypes.Bool, Value: value}

	case int:
		return &GObject{Type: gtypes.Number, Value: float64(value)}

	case float64:
		return &GObject{Type: gtypes.Number, Value: value}

	case string:
		value = strings.ReplaceAll(value, "\\n", "\n")
		return &GObject{Type: gtypes.String, Value: value}

	case func(*GFuncCtx) *GObject:
		return &GObject{Type: gtypes.Func, Value: GFunc(value)}

	}

	assert.Unreached("Could not infer type for %T %v", value, value)
	return nil
}

// Interface returns the *GObject's interface value
func (m *GObject) Interface() interface{} {
	if value, ok := m.Value.(float64); ok && math.Mod(value, 1) == 0 {
		return int(value)
	}

	return m.Value
}

// IsTruthy checks if a *GObject is considered true
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

	case gtypes.Func, gtypes.UserFunc:
		return true
	}

	assert.Unreached("unknown type: %#v", m.Type)
	return false
}

// Call calls a method on a *GObject
func (m *GObject) Call(name string, ctx *GFuncCtx) *GObject {
	assert.True(m.Type.Implements(name), "type %s doesn't implement %s", m.Type.Name, name)

	return m.Type.Resolve(name)(ctx)
}
