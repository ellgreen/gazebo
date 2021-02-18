package g

// ObjectNil is the underlying type of nil in gazebo
type ObjectNil struct {
	object
}

// NewObjectNil creates a new nil object
func NewObjectNil() *ObjectNil {
	return &ObjectNil{object: object{typ: TypeNil}}
}

// Value satisfies the Object interface
func (m *ObjectNil) Value() interface{} {
	return nil
}

// Call satisfies the Object interface
func (m *ObjectNil) Call(method string, args Args) Object {
	return m.call(m, method, args)
}
