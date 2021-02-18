package g

// ObjectInternalFunc is the underlying type of functions in gazebo
type ObjectInternalFunc struct {
	object
	value Func
}

// NewObjectInternalFunc creates a new internal function object
func NewObjectInternalFunc(value Func) *ObjectInternalFunc {
	return &ObjectInternalFunc{
		object: object{typ: TypeInternalFunc},
		value:  value,
	}
}

// Value satisfies the Object interface
func (m *ObjectInternalFunc) Value() interface{} {
	return m.value
}

// Call satisfies the Object interface
func (m *ObjectInternalFunc) Call(method string, args Args) Object {
	return m.call(m, method, args)
}

// Func returns the internal go function
func (m *ObjectInternalFunc) Func() Func {
	return m.value
}
