package g

// ObjectInternalFunc is the underlying type of functions in gazebo
type ObjectInternalFunc struct {
	PartialObject
	value Func
}

// NewObjectInternalFunc creates a new internal function object
func NewObjectInternalFunc(value func(Args) Object) *ObjectInternalFunc {
	return &ObjectInternalFunc{
		PartialObject: PartialObject{typ: TypeInternalFunc},
		value:         Func(value),
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
