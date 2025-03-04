package gsp

import (
	"unsafe"
)

func Zero[T Type]() T {
	switch any(T(0)).(type) {
	case uint8, uint16, uint32, uint64:
		switch unsafe.Sizeof(T(0)) {
		case 1:
			return T(zeroUint8)
		case 2:
			return T(zeroUint16)
		case 4:
			return T(zeroUint32)
		case 8:
			return T(zeroUint64)
		default:
			panic("gsp: Zero: unknown bit size encountered")
		}
	case int8, int16, int32, int64, float32, float64:
		return T(0)
	default:
		panic("gsp: Zero: unsupported type")
	}
}
