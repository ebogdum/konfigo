package validator

// Number is a generic type that can hold either int64 or float64
type Number interface {
	int64 | float64
}

// NumberValue represents a numeric value that can be either integer or float
type NumberValue struct {
	IsFloat bool
	IntVal  int64
	FloatVal float64
}

// ToFloat64 converts the number to float64 for comparison
func (n NumberValue) ToFloat64() float64 {
	if n.IsFloat {
		return n.FloatVal
	}
	return float64(n.IntVal)
}

// FromInterface creates a NumberValue from interface{}
func NumberFromInterface(value interface{}) (NumberValue, bool) {
	switch v := value.(type) {
	case int64:
		return NumberValue{IsFloat: false, IntVal: v}, true
	case float64:
		// Check if it's actually an integer in disguise
		if v == float64(int64(v)) {
			return NumberValue{IsFloat: false, IntVal: int64(v)}, true
		}
		return NumberValue{IsFloat: true, FloatVal: v}, true
	case int:
		return NumberValue{IsFloat: false, IntVal: int64(v)}, true
	case int32:
		return NumberValue{IsFloat: false, IntVal: int64(v)}, true
	case float32:
		return NumberValue{IsFloat: true, FloatVal: float64(v)}, true
	default:
		return NumberValue{}, false
	}
}
