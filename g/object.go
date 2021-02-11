package g

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
)

// New creates a new Object with the provided type and value
func New(typ *Type, value interface{}) Object {
	return &object{typ: typ, value: value}
}

// NewObject creates a new Object for the provided Go value
func NewObject(value interface{}) Object {
	switch value := value.(type) {
	case nil:
		return New(TypeNil, nil)

	case bool:
		return New(TypeBool, value)

	case int:
		return New(TypeNumber, float64(value))

	case float64:
		return New(TypeNumber, value)

	case string:
		return New(TypeString, value)

	case FuncDescription:
		return New(TypeFunc, value)

	case func(Args) Object:
		return New(TypeBuiltinFunc, Func(value))
	}

	assert.Unreached("Could not infer type for %T %v", value, value)
	return nil
}

// NewInternalObject creates a new internal Object
func NewInternalObject(value interface{}) Object {
	return &object{typ: TypeInternal, value: value}
}

// Object is the type of all values in gazebo
type Object interface {
	Type() *Type
	Value() interface{}
	Call(string, Args) Object
	Attributes() *Attributes
}

// FuncDescription contains the code necessary to call
// a user-defined function
type FuncDescription struct {
	Params []string
	Body   compiler.Code
	Env    interface{} // FIXME: create an interface for the VM's env
}

// Attributes is a mapping of attribute names to Object values
type Attributes struct {
	values map[string]Object
}

func (m *Attributes) init() {
	if m.values == nil {
		m.values = make(map[string]Object)
	}
}

// Has returns whether an attribute exists
func (m *Attributes) Has(name string) bool {
	m.init()

	_, ok := m.values[name]
	return ok
}

// Get returns an attribute
func (m *Attributes) Get(name string) Object {
	m.init()

	if m.Has(name) {
		return m.values[name]
	}

	return NewObject(nil)
}

// Set sets and attribute's value
func (m *Attributes) Set(name string, value Object) {
	m.init()

	m.values[name] = value
}

// Delete deletes an attribute
func (m *Attributes) Delete(name string) {
	m.init()

	delete(m.values, name)
}

type object struct {
	typ        *Type
	value      interface{}
	attributes Attributes
}

// Type returns the object's underlying Type
func (m *object) Type() *Type {
	return m.typ
}

// Value returns the object's internal value
func (m *object) Value() interface{} {
	return m.value
}

// Call calls a method on an object
func (m *object) Call(name string, args Args) Object {
	assert.True(m.typ.Implements(name), "type %q does not implement %q", m.typ.Name, name)

	return m.typ.Resolve(name)(m, args)
}

// Attributes returns the Object's attributes
func (m *object) Attributes() *Attributes {
	return &m.attributes
}
