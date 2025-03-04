package gsp

import (
	"math"
	"unsafe"
)

var (
	maxInt8  = math.MaxInt8
	maxInt16 = math.MaxInt16
	maxInt32 = math.MaxInt32
	maxInt64 = math.MaxInt64
)

var (
	zeroUint8  uint8  = -math.MinInt8
	zeroUint16 uint16 = -math.MinInt16
	zeroUint32 uint32 = -math.MinInt32
	zeroUint64 uint64 = -math.MinInt64
)

func isUnsigned[T Type]() bool {
	return T(0)-1 > 0
}

func isSigned[T Type]() bool {
	switch unsafe.Sizeof(T(0)) {
	case 4: // 32-bit
		return T(maxInt32)+1 < 0
	case 8: // 64-bit
		return T(maxInt64)+1 < 0
	default:
		return false
	}
}
