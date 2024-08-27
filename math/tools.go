package math

func IsNumber(arg any) bool {
	return IsFloat(arg) || IsComplex(arg) || IsInteger(arg)
}

func IsInteger(arg any) bool {
	return IsSigned(arg) || IsUnSigned(arg)
}

func IsFloat(arg any) bool {
	switch arg.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

func IsComplex(arg any) bool {
	switch arg.(type) {
	case complex64, complex128:
		return true
	default:
		return false
	}
}

func IsSigned(arg any) bool {
	switch arg.(type) {
	case int, int8, int16, int32, int64:
		return true
	case float32, float64:
		return false
	default:
		return false
	}
}

func IsUnSigned(arg any) bool {
	switch arg.(type) {
	case int, int8, int16, int32, int64:
		return false
	case float32, float64:
		return false
	default:
		return true
	}
}
