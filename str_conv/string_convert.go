package str_conv

import (
	"encoding"
	"reflect"
	"strconv"
	"unsafe"
)

func ToInt8(v string) int8 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseInt(v, 10, 8)
	return int8(x)
}

func ToUint8(v string) uint8 {
	x, _ := strconv.ParseUint(v, 10, 8)
	return uint8(x)
}

func ToInt16(v string) int16 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseInt(v, 10, 16)
	return int16(x)
}

func ToUint16(v string) uint16 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseUint(v, 10, 16)
	return uint16(x)
}

func ToInt32(v string) int32 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseInt(v, 10, 32)
	return int32(x)
}

func ToUint32(v string) uint32 {
	x, _ := strconv.ParseUint(v, 10, 32)
	return uint32(x)
}

func ToInt64(v string) int64 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseInt(v, 10, 64)
	return int64(x)
}

func ToUint64(v string) uint64 {
	x, _ := strconv.ParseUint(v, 10, 64)
	return x
}

func ToInt(v string) int {
	if v == "" {
		return 0
	}
	x, _ := strconv.Atoi(v)
	return x
}

func ToUint(v string) uint {
	x, _ := strconv.ParseUint(v, 10, 0)
	return uint(x)
}

func ToFloat32(v string) float32 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseFloat(v, 32)
	return float32(x)
}

func ToFloat64(v string) float64 {
	if v == "" {
		return 0
	}
	x, _ := strconv.ParseFloat(v, 64)
	return x
}

func ToBool(v string) bool {
	if v == "" {
		return false
	}
	val, _ := strconv.ParseBool(v)
	return val
}

func ToText(v string, dst encoding.TextUnmarshaler) error {
	return dst.UnmarshalText(ToBytes(v))
}

func ToBinary(v string, dst encoding.BinaryUnmarshaler) error {
	return dst.UnmarshalBinary(ToBytes(v))
}

func ToBytes(v string) (b []byte) {
	p := unsafe.StringData(v)
	b = unsafe.Slice(p, len(v))
	return b
}

func ToString(arg any) string {
	switch v := arg.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case encoding.TextMarshaler:
		dst, err := v.MarshalText()
		if err != nil {
			panic(err)
		}
		return unsafe.String(unsafe.SliceData(dst), len(dst))
		//return *(*string)(unsafe.Pointer(&dst))
	case encoding.BinaryMarshaler:
		dst, err := v.MarshalBinary()
		if err != nil {
			panic(err)
		}
		return unsafe.String(unsafe.SliceData(dst), len(dst))
		//return *(*string)(unsafe.Pointer(&dst))
	case string:
		return v
	case []byte:
		if v==nil ||v[0] == 0 {
			return ""
		}
		_max := len(v)
		_min := 0
		n := _max / 2
		for {
			if v[n] == 0 {
				_max = n
				n -= (_max - _min) / 2
				if n == _max {
					break
				}
			} else {
				_min = n
				n += (_max - _min) / 2
				if n == _min {
					break
				}
			}
		}
		return unsafe.String(unsafe.SliceData(v), _max)
	default:
		tv := reflect.ValueOf(arg)
		if tv.Kind() == reflect.Array {
			switch tv.Type().Elem().Kind() {
			case reflect.Uint8:
				if tv.Len() == 0 {
					return ""
				}
				builder := make([]byte, 0, tv.Len())
				for i := 0; i < tv.Len(); i++ {
					x := tv.Index(i).Interface().(byte)
					if x == 0 {
						break
					}
					builder = append(builder, tv.Index(i).Interface().(byte))
				}
				if len(builder) == 0 {
					return ""
				}
				return ToString(builder)
			default:
			}
		}
		panic(tv.Kind().String() + " is not supported when converting to string")
	}
}
