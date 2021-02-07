package assert

import "github.com/kr/pretty"

// Error is panicked with when an assertion fails
type Error struct {
	message string
}

// Error satisfies the error interface
func (e Error) Error() string {
	return e.message
}

func fail(args ...interface{}) {
	message := "Assertion failed"

	if len(args) > 0 {
		message += ": " + pretty.Sprintf(args[0].(string), args[1:]...)
	}

	panic(Error{message: message})
}

// True asserts that a condition is true
func True(condition bool, args ...interface{}) {
	if !condition {
		fail(args...)
	}
}

// False asserts that a condition is false
func False(condition bool, args ...interface{}) {
	if condition {
		fail(args...)
	}
}

// Nil asserts that a value is nil
func Nil(value interface{}, args ...interface{}) {
	if value != nil {
		fail(args...)
	}
}

// NotNil asserts a value is not nil
func NotNil(value interface{}, args ...interface{}) {
	if value == nil {
		fail(args...)
	}
}

// Unreached always results in a failed assertion
func Unreached(args ...interface{}) {
	fail(args...)
}
