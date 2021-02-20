package g

// ObjectBool is the underlying type of booleans in gazebo
type ObjectBool struct {
	PartialObject
	value bool
}

// NewObjectBool creates a new boolean object
func NewObjectBool(value bool) *ObjectBool {
	return &ObjectBool{
		PartialObject: PartialObject{typ: TypeBool},
		value:         value,
	}
}

// Value satisfies the Object interface
func (m *ObjectBool) Value() interface{} {
	return m.value
}

// Call satisfies the Object interface
func (m *ObjectBool) Call(method string, args Args) Object {
	return m.call(m, method, args)
}

// Bool returns the boolean value of an ObjectBool
func (m *ObjectBool) Bool() bool {
	return m.value
}
