package gsp

import (
	"math"
	"unsafe"
)

var (
	maxInt8   int8   = math.MaxInt8
	maxUint8  uint8  = math.MaxUint8
	maxInt16  int16  = math.MaxInt16
	maxUint16 uint16 = math.MaxUint16
	maxInt32  int32  = math.MaxInt32
	maxUint32 uint32 = math.MaxUint32
)

var (
	minInt8    int8    = -math.MaxInt8
	minInt16   int16   = -math.MaxInt16
	minInt32   int32   = -math.MaxInt32
	minFloat32 float32 = -1
	minFloat64 float64 = -1
)

const (
	maxInt8_32  float32 = math.MaxInt8
	maxInt16_32 float32 = math.MaxInt16
	maxInt32_32 float32 = math.MaxInt32

	invMaxInt8_32  float32 = 1.0 / maxInt8_32
	invMaxInt16_32 float32 = 1.0 / maxInt16_32
	invMaxInt32_32 float32 = 1.0 / maxInt32_32
)

const (
	maxInt8_64  float64 = math.MaxInt8
	maxInt16_64 float64 = math.MaxInt16
	maxInt32_64 float64 = math.MaxInt32

	invMaxInt8_64  float64 = 1.0 / maxInt8_64
	invMaxInt16_64 float64 = 1.0 / maxInt16_64
	invMaxInt32_64 float64 = 1.0 / maxInt32_64
)

var (
	zeroUint8  uint8  = -math.MinInt8
	zeroUint16 uint16 = -math.MinInt16
	zeroUint32 uint32 = -math.MinInt32
)

func isUnsigned[T Type]() bool {
	return T(0)-1 > 0
}

func isSigned[T Type]() bool {
	switch unsafe.Sizeof(T(0)) {
	case 1: // 8-bit
		return T(maxInt8)+1 < 0
	case 2: // 16-bit
		return T(maxInt16)+1 < 0
	case 4: // 32-bit
		return T(maxInt32)+1 < 0
	default:
		return false
	}
}
