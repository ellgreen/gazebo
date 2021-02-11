package g

// IsTruthy determines if the provided Object is truthy
func IsTruthy(object Object) bool {
	return object.Call(Protocols.ToBool, nil).Value().(bool)
}

// ToString returns an Object's string value
func ToString(object Object) string {
	return object.Call(Protocols.ToString, nil).Value().(string)
}

// ToNumber returns an Object's numeric value
func ToNumber(object Object) float64 {
	return object.Call(Protocols.ToNumber, nil).Value().(float64)
}
