package math

import (
	"unsafe"
)

type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnSignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Integer interface {
	SignedInteger | UnSignedInteger
}

type Number interface {
	Integer | ~float32 | ~float64
}

// CeilToPowerOfTwo 将整数型的数扩大到2^n,返回0代表执行错误
func CeilToPowerOfTwo[T Integer](n T) T {
	if n <= 0 {
		return 0
	}
	if n <= 2 {
		return n
	}
	var offset = []int{1, 2, 4, 8, 16, 32}
	bit := uint32(unsafe.Sizeof(T(0))) //计算T的字节数

	times := 3
	switch bit {
	case 16:
		times = 4
	case 32:
		times = 5
	case 64:
		times = 6
	}

	n--
	for i := 0; i < times; i++ {
		n |= n >> offset[i]
	}
	n++
	if n < 0 {
		n = 0
	}
	return n
}

func FloorToPowerOfTwo[T Integer](n T) T {
	if n <= 0 {
		return 0
	}
	if n <= 2 {
		return n
	}
	var bit = uint32(unsafe.Sizeof(T(0))) //计算T的字节数
	var offset = []int{1, 2, 4, 8, 16, 32}
	times := 3
	switch bit {
	case 16:
		times = 4
	case 32:
		times = 5
	case 64:
		times = 6
	}
	for i := 0; i < times; i++ {
		n |= n >> offset[i]
	}
	n >>= 1
	n++
	return n
}
