package g

import (
	"github.com/johnfrankmorgan/gazebo/assert"
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

// partial Object implementation
type object struct {
	typ        *Type
	attributes Attributes
}

// Type returns the object's underlying Type
func (m *object) Type() *Type {
	return m.typ
}

// Attributes returns the Object's attributes
func (m *object) Attributes() *Attributes {
	return &m.attributes
}

func (m *object) call(self Object, method string, args Args) Object {
	assert.True(m.typ.Implements(method), "type %q does not implement %q", m.typ.Name, method)

	return m.typ.Resolve(method)(self, args)
}
