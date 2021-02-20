package g

// ObjectNumber is the underlying type of numbers in gazebo
type ObjectNumber struct {
	PartialObject
	value float64
}

// NewObjectNumber creates a new number object
func NewObjectNumber(value float64) *ObjectNumber {
	return &ObjectNumber{
		PartialObject: PartialObject{typ: TypeNumber},
		value:         value,
	}
}

// Value satisfies the Object interface
func (m *ObjectNumber) Value() interface{} {
	return m.value
}

// Call satisfies the Object interface
func (m *ObjectNumber) Call(method string, args Args) Object {
	return m.call(m, method, args)
}

// Float returns the float value of an ObjectNumber
func (m *ObjectNumber) Float() float64 {
	return m.value
}

// Int returns the integer value of an ObjectNumber
func (m *ObjectNumber) Int() int {
	return int(m.value)
}
