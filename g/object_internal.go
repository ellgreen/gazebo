package g

import "github.com/johnfrankmorgan/gazebo/assert"

// ObjectInternal is used for values internal to the gazebo VM
type ObjectInternal struct {
	object
	value interface{}
}

// NewObjectInternal creates a new internal object
func NewObjectInternal(value interface{}) *ObjectInternal {
	return &ObjectInternal{
		object: object{typ: TypeInternal},
		value:  value,
	}
}

// Value satisfies the Object interface
func (m *ObjectInternal) Value() interface{} {
	return m.value
}

// Call satisfies the Object interface
func (m *ObjectInternal) Call(method string, args Args) Object {
	assert.Unreached()
	return nil
}
