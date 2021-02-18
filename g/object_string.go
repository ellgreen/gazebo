package g

// ObjectString is the underlying type of strings in gazebo
type ObjectString struct {
	object
	value string
}

// NewObjectString creates a new string object
func NewObjectString(value string) *ObjectString {
	return &ObjectString{
		object: object{typ: TypeString},
		value:  value,
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
