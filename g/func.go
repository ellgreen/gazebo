package g

import "github.com/johnfrankmorgan/gazebo/assert"

// Func is the type of gazebo's builtin functions
type Func func(args Args) Object

// Args is used to pass arguments to builtin gazebo functions
type Args []Object

// Expects asserts that N arguments were received
func (m Args) Expects(n int) {
	assert.Len(m, n, "expected %d arguments, got %d", n, len(m))
}

// ExpectsAtLeast asserts that at least N arguments were received
func (m Args) ExpectsAtLeast(n int) {
	assert.True(len(m) >= n, "expected at least %d arguments, got %d", n, len(m))
}

// Self is a helper method to return the first argument
func (m Args) Self() Object {
	m.ExpectsAtLeast(1)
	return m[0]
}

// SelfWithArgs returns the first argument as an Object
// and the remaining args in an Args slice
func (m Args) SelfWithArgs() (Object, Args) {
	return m.Self(), m[1:]
}

// Values returns an array of interface{} values from an
// Args slice
func (m Args) Values() []interface{} {
	var ifaces []interface{}

	for _, arg := range m {
		ifaces = append(ifaces, arg.Value())
	}

	return ifaces
}

// Parse parses Object values into the specified go pointers
func (m Args) Parse(dests ...interface{}) {
	m.ExpectsAtLeast(len(dests))

	for i, dest := range dests {
		value := m[i].Value()

		switch dest := dest.(type) {
		case *bool:
			*dest = value.(bool)

		case *int:
			*dest = int(value.(float64))

		case *float64:
			*dest = value.(float64)

		case *string:
			*dest = value.(string)

		case *Object:
			*dest = m[i]

		default:
			assert.Unreached("cannot parse arg type %T", dest)
		}
	}
}
