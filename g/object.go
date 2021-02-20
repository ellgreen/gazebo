package g

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/errors"
)

// Object is the type of all values in gazebo
type Object interface {
	Type() *Type
	Value() interface{}
	Call(string, Args) Object
	Attributes() *Attributes
}

// NewObject creates a new Object for the provided Go value
func NewObject(value interface{}) Object {
	switch value := value.(type) {
	case nil:
		return NewObjectNil()

	case bool:
		return NewObjectBool(value)

	case int:
		return NewObjectNumber(float64(value))

	case float64:
		return NewObjectNumber(value)

	case string:
		return NewObjectString(value)

	case []Object:
		return NewObjectList(value)

	case func(Args) Object:
		return NewObjectInternalFunc(Func(value))
	}

	assert.Unreached("Could not infer type for %T %v", value, value)
	return nil
}

// PartialObject is a partial Object implementation
type PartialObject struct {
	typ        *Type
	attributes Attributes
}

// Type returns the object's underlying Type
func (m *PartialObject) Type() *Type {
	return m.typ
}

// Attributes returns the Object's attributes
func (m *PartialObject) Attributes() *Attributes {
	return &m.attributes
}

func (m *PartialObject) call(self Object, method string, args Args) Object {
	errors.ErrRuntime.Expect(
		m.typ.Implements(method),
		"type %s does not implement %s",
		m.typ.Name,
		method,
	)

	return m.typ.Resolve(method)(self, args)
}

// CallMethod is an exported method that wraps PartialObject.call
func (m *PartialObject) CallMethod(self Object, method string, args Args) Object {
	return m.call(self, method, args)
}

// SetType is an exported method to allow modules to set an Object's type
func (m *PartialObject) SetType(typ *Type) {
	m.typ = typ
}
