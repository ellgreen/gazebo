package g

// ObjectString is the underlying type of strings in gazebo
type ObjectString struct {
	PartialObject
	value string
}

// NewObjectString creates a new string object
func NewObjectString(value string) *ObjectString {
	return &ObjectString{
		PartialObject: PartialObject{typ: TypeString},
		value:         value,
	}
}

// Value satisfies the Object interface
func (m *ObjectString) Value() interface{} {
	return m.value
}

// Call satisfies the Object interface
func (m *ObjectString) Call(method string, args Args) Object {
	return m.call(m, method, args)
}

// String returns the object's string value
func (m *ObjectString) String() string {
	return m.value
}

// Len returns the object's length
func (m *ObjectString) Len() int {
	return len(m.value)
}
